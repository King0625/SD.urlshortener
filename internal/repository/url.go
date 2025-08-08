package repository

import (
	"context"

	"github.com/King0625/SD.urlshortener/internal/db/sqlc"
	"github.com/jackc/pgx/v5"
)

type UrlRepository interface {
	Create(ctx context.Context, url sqlc.CreateURLParams) (sqlc.Url, error)
	GetOneByCode(ctx context.Context, code string) (sqlc.Url, error)
}

type urlRepository struct {
	Queries *sqlc.Queries
	Conn    *pgx.Conn
}

func NewUrlRepository(conn *pgx.Conn) UrlRepository {
	queries := sqlc.New(conn)
	return &urlRepository{
		Queries: queries,
		Conn:    conn,
	}
}

func (r *urlRepository) Create(ctx context.Context, url sqlc.CreateURLParams) (sqlc.Url, error) {
	return r.Queries.CreateURL(ctx, url)
}

func (r *urlRepository) GetOneByCode(ctx context.Context, code string) (sqlc.Url, error) {
	return r.Queries.GetURLByCode(ctx, code)
}
