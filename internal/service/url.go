package service

import (
	"context"
	"math/rand"

	"github.com/King0625/SD.urlshortener/internal/db/sqlc"
	"github.com/King0625/SD.urlshortener/internal/repository"
)

const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateCode(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(b)
}

type UrlService interface {
	ShortenUrl(ctx context.Context, originalUrl string) (sqlc.Url, error)
	GetUrlByCode(ctx context.Context, code string) (sqlc.Url, error)
}

type urlService struct {
	repo repository.UrlRepository
}

func NewUrlService(repo repository.UrlRepository) UrlService {
	return &urlService{repo: repo}
}

func (s *urlService) ShortenUrl(ctx context.Context, originalUrl string) (sqlc.Url, error) {
	code := generateCode(6)
	todo := sqlc.CreateURLParams{
		Code:        code,
		OriginalUrl: originalUrl,
	}
	return s.repo.Create(ctx, todo)
}

func (s *urlService) GetUrlByCode(ctx context.Context, code string) (sqlc.Url, error) {
	return s.repo.GetOneByCode(ctx, code)
}
