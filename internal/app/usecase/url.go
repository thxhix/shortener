package usecase

import (
	"github.com/thxhix/shortener/internal/app/database"
)

type URLUseCase struct {
	database database.Database
}

func GetInstance() *URLUseCase {
	return &URLUseCase{database: *database.CreateDatabase()}
}

func (u *URLUseCase) Shorten(originalURL string) (string, error) {
	shorten := GetHash()
	shorten, err := u.database.AddLink(originalURL, shorten)
	if err != nil {
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
