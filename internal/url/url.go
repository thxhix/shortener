package url

import (
	"context"
	"errors"
	"github.com/thxhix/shortener/internal/config"
	"github.com/thxhix/shortener/internal/database/interfaces"
	customErrors "github.com/thxhix/shortener/internal/errors"
	"github.com/thxhix/shortener/internal/middleware"
	"github.com/thxhix/shortener/internal/models"
	"sync"
)

type URLUseCaseInterface interface {
	Shorten(ctx context.Context, url string) (string, error)
	GetFullURL(url string) (string, error)
	PingDB() error
	BatchShorten(ctx context.Context, list models.BatchShortenRequestList) (models.BatchShortenResponseList, error)
	UserList(userID string) (models.UserLinksResponseList, error)
	UserDeleteRows(userID string, ids []string) error
}

var ErrLinkDeleted = errors.New("DELETED")

type URLUseCase struct {
	database interfaces.Database
	cfg      *config.Config
}

func NewURLUseCase(db interfaces.Database, cfg config.Config) *URLUseCase {
	return &URLUseCase{
		database: db,
		cfg:      &cfg,
	}
}

func (u *URLUseCase) Shorten(ctx context.Context, originalURL string) (string, error) {
	shorten := GetHash()
	shorten, err := u.database.AddLink(originalURL, shorten, middleware.GetUserID(ctx))
	if err != nil {
		if errors.Is(err, customErrors.ErrDuplicate) {
			return shorten, customErrors.ErrDuplicate
		}
		return "", err
	}
	return shorten, nil
}

func (u *URLUseCase) GetFullURL(hash string) (string, error) {
	link, err := u.database.GetFullLink(hash)
	if err != nil {
		return "", err
	}
	if link.IsDeleted {
		return "", ErrLinkDeleted
	}
	return link.URL, nil
}

func (u *URLUseCase) PingDB() error {
	return u.database.PingConnection()
}

func (u *URLUseCase) BatchShorten(ctx context.Context, list models.BatchShortenRequestList) (models.BatchShortenResponseList, error) {
	var result models.DBShortenRowList
	var response models.BatchShortenResponseList

	// Будто очень кривое исполнение, но надо сдать спринт..
	// TODO: глянуть, отрефакторить
	for _, batch := range list {
		row := models.DBShortenRow{
			Hash: GetHash(),
			URL:  batch.URL,
		}
		result = append(result, row)

		responseRow := models.BatchShortenResponse{
			ID:   batch.ID,
			Hash: u.cfg.BaseURL + "/" + row.Hash,
		}
		response = append(response, responseRow)
	}

	err := u.database.AddLinks(ctx, result, middleware.GetUserID(ctx))
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (u *URLUseCase) UserList(userID string) (models.UserLinksResponseList, error) {
	links, err := u.database.GetUserFullLinks(userID)
	if err != nil {
		return nil, err
	}

	result := models.UserLinksResponseList{}

	for _, link := range links {
		row := models.UserLinksResponse{
			Original: link.URL,
			Short:    u.cfg.BaseURL + "/" + link.Hash,
		}
		result = append(result, row)
	}

	return result, nil
}

func (u *URLUseCase) UserDeleteRows(userID string, ids []string) error {
	numWorkers := 10
	batchSize := 1000

	var wg sync.WaitGroup
	batchCh := make(chan []string)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for batch := range batchCh {
				_ = u.database.RemoveUserLinks(userID, batch)
			}
		}()
	}

	for i := 0; i < len(ids); i += batchSize {
		end := i + batchSize
		if end > len(ids) {
			end = len(ids)
		}
		batch := ids[i:end]
		batchCh <- batch
	}

	close(batchCh)
	wg.Wait()

	return nil
}
