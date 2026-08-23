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

func (m *MockLinkRepository) Create(ctx context.Context, link *models.Link) error {
	args := m.Called(ctx, link)
	return args.Error(0)
}

func (m *MockLinkRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockLinkRepository) GetAll(ctx context.Context) ([]models.Link, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Link), args.Error(1)

}

func (m *MockLinkRepository) GetByShorten(ctx context.Context, shortLink string) (string, error) {
	args := m.Called(ctx, shortLink)

	if args.Get(0) == nil {
		return "", args.Error(1)
	}

	return args.Get(0).(string), args.Error(1)
}

func TestCreateLink_Success(t *testing.T) {
	mockRepo := new(MockLinkRepository)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	service := services.NewLinkService(mockRepo, logger)

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
	service := services.NewLinkService(mockRepo, nil)

	link, err := service.CreateLink(context.Background(), "not a url")

	assert.ErrorIs(t, err, apperrors.ErrInvalidURL)
	assert.Nil(t, link)

	mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)

}

func TestCreateLink_DBError(t *testing.T) {
	mockRepo := new(MockLinkRepository)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	service := services.NewLinkService(mockRepo, logger)

	dbErr := errors.New("database connection lost")

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Link")).Return(dbErr)

	link, err := service.CreateLink(context.Background(), "https://www.google.com")

	assert.Error(t, err)
	assert.Nil(t, link)
	assert.Contains(t, err.Error(), "connection")

}
