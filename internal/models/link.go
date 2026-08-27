package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Link struct {
	ID        int
	Origin    string `redis:"origin"`
	Shorten   string `redis:"shorten"`
	CreatedAt time.Time
}

type LinkModel struct {
	DB *pgxpool.Pool
}
