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

var ()

type LinkService struct {
	repo   repository.LinkRepository
	logger *slog.Logger
}

func NewLinkService(r repository.LinkRepository, logger *slog.Logger) *LinkService {
	return &LinkService{
		repo:   r,
		logger: logger,
	}
}

func (l *LinkService) CreateLink(ctx context.Context, originUrl string) (*models.Link, error) {
	if !isValidUrl(originUrl) {
		return nil, apperrors.ErrInvalidURL
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

	l.logger.Info("Link created successfuly", "shorten", link.Shorten)
	return link, nil

}

func (l *LinkService) DeleteLink(ctx context.Context, id int) error {
	if err := l.repo.Delete(ctx, id); err != nil {
		l.logger.Error("Failed to delete link by id ")
		return fmt.Errorf("Failed to delete link by id %v: %w", id, err)
	}

	return nil
}

func (l *LinkService) GetAllLinks(ctx context.Context) ([]models.Link, error) {
	links, err := l.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error in gettin all links")
	}

	l.logger.Info("All links got successfuly")
	return links, nil
}

func (l *LinkService) GetByShorten(ctx context.Context, short string) (string, error) {
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
