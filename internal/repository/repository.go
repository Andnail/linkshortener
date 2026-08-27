package repository

import (
	"context"
	"linkshortener/internal/models"
)

type LinkRepository interface {
	//GetAll(ctx context.Context) ([]models.Link, error)
	GetByOrigin(ctx context.Context, originLink string) (string, error)
	GetByShorten(ctx context.Context, shortLink string) (string, error)
	Delete(ctx context.Context, url string) error
	Create(ctx context.Context, link *models.Link) error
}

type CacheRepository interface {
	Set(ctx context.Context, link *models.Link) error
	Delete(ctx context.Context, origin string) error
	GetByString(ctx context.Context, origin string) (string, error)
}
