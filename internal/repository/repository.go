package repository

import (
	"context"
	"linkshortener/internal/models"
)

type LinkRepository interface {
	GetAll(ctx context.Context) ([]models.Link, error)
	GetByShorten(ctx context.Context, shortLink string) (string, error)
	Delete(ctx context.Context, id int) error
	Create(ctx context.Context, link *models.Link) error
}
