package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"linkshortener/internal/apperrors"
	"linkshortener/internal/models"
	"linkshortener/internal/repository"
	"log/slog"
	"net/url"
	"time"
)

type LinkService struct {
	repo   repository.LinkRepository
	cache  repository.CacheRepository
	logger *slog.Logger
}

func NewLinkService(r repository.LinkRepository, cache repository.CacheRepository, logger *slog.Logger) *LinkService {
	return &LinkService{
		repo:   r,
		cache:  cache,
		logger: logger,
	}
}

func (l *LinkService) CreateLink(ctx context.Context, originUrl string) (*models.Link, error) {
	if !isValidUrl(originUrl) {
		return nil, apperrors.ErrInvalidURL
	}

	if link, err := l.cache.GetByString(ctx, originUrl); err == nil && link != "" {
		l.logger.Info("Found in cache")
		return &models.Link{
			Origin:  originUrl,
			Shorten: link,
		}, nil
	}

	link := &models.Link{
		Origin:    originUrl,
		Shorten:   generateRandomString(),
		CreatedAt: time.Now(),
	}

	if err := l.repo.Create(ctx, link); err != nil {
		l.logger.Error("failed to create a link", "error", err, "url", originUrl)
		return nil, fmt.Errorf("Create link %w", err)
	}

	if err := l.cache.Set(ctx, link); err != nil {
		l.logger.Error("Failed to add link to cache", "error", err)
	}

	l.logger.Info("Link created successfuly", "shorten", link.Shorten)
	return link, nil

}

func (l *LinkService) DeleteLink(ctx context.Context, url string) error {
	if err := l.repo.Delete(ctx, url); err != nil {
		l.logger.Error("Failed to delete link by origin url ")
		return fmt.Errorf("Failed to delete link by origin url %v: %w", url, err)
	}

	return nil
}

func (l *LinkService) GetByShorten(ctx context.Context, short string) (string, error) {
	if link, err := l.cache.GetByString(ctx, short); err != nil {
		return link, nil
	}

	origin, err := l.repo.GetByShorten(ctx, short)
	if err != nil {
		return "", apperrors.ErrLinkNotFound
	}

	return origin, nil
}

func isValidUrl(link string) bool {
	_, err := url.ParseRequestURI(link)
	return err == nil
}

func generateRandomString() string {
	length := 10
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		panic("failed to generate random string")
	}

	return string(hex.EncodeToString(bytes)[:length])
}
