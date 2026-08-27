package postgresrepo

import (
	//"linkshortener/internal/logger"
	"context"
	"fmt"

	"linkshortener/internal/models"
	"linkshortener/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) repository.LinkRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (p *PostgresRepository) Create(ctx context.Context, link *models.Link) error {
	stmt := `INSERT INTO links (origin, shorten, created_at)
	VALUES ( $1, $2,$3);`

	_, err := p.db.Exec(ctx, stmt,
		link.Origin,
		link.Shorten,
		link.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("Creating link: %w", err)
	}

	return nil
}

func (p *PostgresRepository) Delete(ctx context.Context, url string) error {
	stmt := `DELETE FROM links WHERE origin = $1`

	_, err := p.db.Exec(ctx, stmt, url)
	if err != nil {
		return fmt.Errorf("Delete link: %w", err)
	}

	return nil
}

// func (p *PostgresRepository) GetAll(ctx context.Context) ([]models.Link, error) {
// 	stmt := `SELECT id, origin, shorten, created_at
// 		FROM links
// 		ORDER BY created_at DESC`
// 	rows, err := p.db.Query(ctx, stmt)
// 	if err != nil {
// 		return nil, fmt.Errorf("Query links %w", err)
// 	}
// 	defer rows.Close()

// 	allLinks := make([]models.Link, 0)

// 	for rows.Next() {
// 		var l models.Link
// 		err := rows.Scan(&l.ID, &l.Origin, &l.Shorten, &l.CreatedAt)
// 		if err != nil {
// 			return nil, fmt.Errorf("links scan: %w", err)
// 		}
// 		allLinks = append(allLinks, l)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("iterate links: %w", err)
// 	}

// 	return allLinks, nil
// }

func (p *PostgresRepository) GetByOrigin(ctx context.Context, originLink string) (string, error) {
	var l string
	stmt := `SELECT id, origin, shorten, created_at
 		FROM links WHERE origin = $1`

	rows := p.db.QueryRow(ctx, stmt, originLink)

	err := rows.Scan(&l)
	if err != nil {
		return "", fmt.Errorf("link scan: %w", err)
	}

	return l, nil
}

func (p *PostgresRepository) GetByShorten(ctx context.Context, shortLink string) (string, error) {
	var originLink string
	stmt := `SELECT origin FROM links WHERE shorten = $1`

	err := p.db.QueryRow(ctx, stmt, shortLink).Scan(&originLink)
	if err != nil {
		return "", fmt.Errorf("Failed to get by short link")
	}

	return originLink, nil
}

// func GetByShortCode(ctx context.Context, db *pgxpool.Pool, code string) (*models.Link, error){
// 	stmt := `SELECT * FROM links WHERE short_code`
// }
