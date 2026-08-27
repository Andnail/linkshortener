package services_test

import (
	"context"
	"errors"
	"io"
	"linkshortener/internal/apperrors"
	"linkshortener/internal/models"
	"linkshortener/internal/services"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLinkRepository struct {
	mock.Mock
}

type MockCacheRepository struct {
	mock.Mock
}

func (l *MockLinkRepository) Create(ctx context.Context, link *models.Link) error {
	args := l.Called(ctx, link)
	return args.Error(0)
}

func (l *MockLinkRepository) Delete(ctx context.Context, url string) error {
	args := l.Called(ctx, url)
	return args.Error(0)
}

// func (m *MockLinkRepository) GetAll(ctx context.Context) ([]models.Link, error) {
// 	args := m.Called(ctx)
// 	return args.Get(0).([]models.Link), args.Error(1)

// }

func (l *MockLinkRepository) GetByShorten(ctx context.Context, shortLink string) (string, error) {
	args := l.Called(ctx, shortLink)

	if args.Get(0) == nil {
		return "", args.Error(1)
	}

	return args.Get(0).(string), args.Error(1)
}

func (l *MockLinkRepository) GetByOrigin(ctx context.Context, originLink string) (string, error) {
	return "", nil
}

func (c *MockCacheRepository) Set(ctx context.Context, link *models.Link) error {
	args := c.Called(ctx, link)
	return args.Error(0)
}

func (c *MockCacheRepository) Delete(ctx context.Context, url string) error {
	args := c.Called(ctx, url)
	return args.Error(0)
}

func (c *MockCacheRepository) GetByShorten(ctx context.Context, shortLink string) (string, error) {
	return "", nil
}

func (c *MockCacheRepository) GetByString(ctx context.Context, originUrl string) (string, error) {
	args := c.Called(ctx, originUrl)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}

	return args.Get(0).(string), nil
}

func TestCreateLink_Success(t *testing.T) {
	mockRepo := new(MockLinkRepository)
	mockCache := new(MockCacheRepository)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	service := services.NewLinkService(mockRepo, mockCache, logger)

	origin := "https://www.google.com"

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Link")).Return(nil)

	link, err := service.CreateLink(context.Background(), origin)

	assert.NoError(t, err)
	assert.NotNil(t, link)
	assert.Equal(t, origin, link.Origin)
	assert.NotEmpty(t, link.Shorten)

	mockRepo.AssertExpectations(t)
}

func TestCreate_InvalidURL(t *testing.T) {
	mockRepo := new(MockLinkRepository)
	mockCache := new(MockCacheRepository)

	service := services.NewLinkService(mockRepo, mockCache, nil)

	link, err := service.CreateLink(context.Background(), "not a url")

	assert.ErrorIs(t, err, apperrors.ErrInvalidURL)
	assert.Nil(t, link)

	mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)

}

func TestCreateLink_DBError(t *testing.T) {
	mockRepo := new(MockLinkRepository)
	mockCache := new(MockCacheRepository)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	service := services.NewLinkService(mockRepo, mockCache, logger)

	dbErr := errors.New("database connection lost")

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Link")).Return(dbErr)

	link, err := service.CreateLink(context.Background(), "https://www.google.com")

	assert.Error(t, err)
	assert.Nil(t, link)
	assert.Contains(t, err.Error(), "connection")

}
