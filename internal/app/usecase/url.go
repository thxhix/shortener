package usecase

import (
	"context"
	"errors"
	"github.com/thxhix/shortener/internal/app/config"
	"github.com/thxhix/shortener/internal/app/database/interfaces"
	customErrors "github.com/thxhix/shortener/internal/app/errors"
	"github.com/thxhix/shortener/internal/app/models"
)

type URLUseCaseInterface interface {
	Shorten(url string) (string, error)
	GetFullURL(url string) (string, error)
	PingDB() error
	BatchShorten(ctx context.Context, list models.BatchShortenRequestList) (models.BatchShortenResponseList, error)
}

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

func (u *URLUseCase) Shorten(originalURL string) (string, error) {
	shorten := GetHash()
	shorten, err := u.database.AddLink(originalURL, shorten)
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
	return link, nil
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

	err := u.database.AddLinks(ctx, result)
	if err != nil {
		return nil, err
	}

	return response, nil
}
