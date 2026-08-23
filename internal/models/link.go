package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Link struct {
	ID        int
	Origin    string
	Shorten   string
	CreatedAt time.Time
}

type LinkModel struct {
	DB *pgxpool.Pool
}
