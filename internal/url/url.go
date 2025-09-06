package url

import (
	"context"
	"errors"
	"github.com/thxhix/shortener/internal/config"
	"github.com/thxhix/shortener/internal/database/interfaces"
	customErrors "github.com/thxhix/shortener/internal/errors"
	"github.com/thxhix/shortener/internal/middleware"
	"github.com/thxhix/shortener/internal/models"
	"log"
	"sync"
	"time"
)

// URLUseCaseInterface defines the business logic of the URL shortener service.
// This interface is used to work with different storage implementations.
type URLUseCaseInterface interface {
	// Shorten generates a short link for the given URL and saves it to the database.
	// If the URL already exists, returns the same short link with ErrDuplicate.
	Shorten(ctx context.Context, url string) (string, error)

	// GetFullURL returns the original URL by its short hash.
	// If the link has been deleted, returns ErrLinkDeleted.
	GetFullURL(ctx context.Context, hash string) (string, error)

	// PingDB checks the database connection.
	PingDB() error

	// BatchShorten accepts a list of URLs and stores them in the database in batch mode.
	// Returns a list of generated short links.
	BatchShorten(ctx context.Context, list models.BatchShortenRequestList) (models.BatchShortenResponseList, error)

	// UserList returns all links of a user by their userID.
	UserList(ctx context.Context, userID string) (models.UserLinksResponseList, error)

	// UserDeleteRows deletes a set of user links concurrently.
	// numWorkers – number of workers (goroutines).
	// batchSize – number of links processed per worker batch.
	UserDeleteRows(userID string, ids []string, numWorkers int, batchSize int)
}

// ErrLinkDeleted is returned when a deleted link is requested.
var ErrLinkDeleted = errors.New("DELETED")

// URLUseCase is the main implementation of the business logic of the URL shortener service.
// It operates on top of the abstract Database interface.
type URLUseCase struct {
	database interfaces.Database
	cfg      *config.Config
}

// NewURLUseCase creates a new instance of URLUseCase with the given database and config.
func NewURLUseCase(db interfaces.Database, cfg config.Config) *URLUseCase {
	return &URLUseCase{
		database: db,
		cfg:      &cfg,
	}
}

// Shorten generates a short link for the provided original URL and saves it to the database.
// If the link already exists, returns the existing short link with ErrDuplicate.
func (u *URLUseCase) Shorten(ctx context.Context, originalURL string) (string, error) {
	shorten := GetHash()
	shorten, err := u.database.AddLink(ctx, originalURL, shorten, middleware.GetUserID(ctx))
	if err != nil {
		if errors.Is(err, customErrors.ErrDuplicate) {
			return shorten, customErrors.ErrDuplicate
		}
		return "", err
	}
	return shorten, nil
}

// GetFullURL returns the original URL by the given short hash.
// If the link was deleted, returns ErrLinkDeleted.
func (u *URLUseCase) GetFullURL(ctx context.Context, hash string) (string, error) {
	link, err := u.database.GetFullLink(ctx, hash)
	if err != nil {
		return "", err
	}
	if link.IsDeleted {
		return "", ErrLinkDeleted
	}
	return link.URL, nil
}

// PingDB checks if the database connection is alive.
func (u *URLUseCase) PingDB() error {
	return u.database.PingConnection()
}

// BatchShorten accepts a list of URLs and saves them in the database in batch mode.
// Returns a list of generated short links with their IDs.
func (u *URLUseCase) BatchShorten(ctx context.Context, list models.BatchShortenRequestList) (models.BatchShortenResponseList, error) {
	var result models.DBShortenRowList
	var response models.BatchShortenResponseList

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

// UserList returns all user links with full short URLs.
func (u *URLUseCase) UserList(ctx context.Context, userID string) (models.UserLinksResponseList, error) {
	links, err := u.database.GetUserFullLinks(ctx, userID)
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

// UserDeleteRows deletes user links in batches using multiple workers.
// Each worker receives a batch of IDs from the channel and removes them from the database.
// The channel is closed after all batches are sent, and the method waits for all goroutines to finish.
func (u *URLUseCase) UserDeleteRows(userID string, ids []string, numWorkers int, batchSize int) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	batchCh := make(chan []string)

	for i := 0; i < numWorkers; i++ {
		workerID := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			for batch := range batchCh {
				if err := u.database.RemoveUserLinks(ctx, userID, batch); err != nil {
					log.Printf("[worker %d] UserDeleteRows ошибка при удалении ссылок: %v", workerID, err)
					continue
				}
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
}
