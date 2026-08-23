package apperrors

import (
	"errors"
)

// обработка неудачного подключения к Postgress
var (
	connectionPostgresError = errors.New("Failed to connect to Postgres")
	insertPostgresError     = errors.New("Failed to insert data into table")
	ErrInvalidURL           = errors.New("Invalid URL")
	ErrLinkNotFound         = errors.New("Link not found")
)
