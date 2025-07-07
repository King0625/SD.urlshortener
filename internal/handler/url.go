package handler

import (
	"math/rand"
	"net/http"

	"github.com/King0625/SD.urlshortener/internal/db"
	"github.com/King0625/SD.urlshortener/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateCode(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(b)
}

type UrlHandler struct {
	Queries *db.Queries
	Conn    *pgx.Conn
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"shorturl"`
}

func (h *UrlHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var message string
	var req ShortenRequest
	if err := utils.ReadJSONRequest(w, r, &req); err != nil || req.URL == "" {
		message = "invalid request"
		utils.RespondError(w, http.StatusBadRequest, "BAD_REQUEST", message, nil)
		return
	}

	code := generateCode(6)

	result, err := h.Queries.CreateURL(r.Context(), db.CreateURLParams{
		Code:        code,
		OriginalUrl: req.URL,
	})

	if err != nil {
		message = "create short url failed"
		utils.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", message, nil)
		return
	}

	res := ShortenResponse{
		ShortURL: "http://localhost:9527/" + result.Code,
	}

	message = "Create short url successfully"
	utils.RespondSuccess(w, http.StatusCreated, message, res)
}

func (h *UrlHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	var message string
	code := chi.URLParam(r, "code")

	urlRow, err := h.Queries.GetURLByCode(r.Context(), code)
	if err != nil {
		message = "short url not found"
		utils.RespondError(w, http.StatusNotFound, "SHORT_URL_NOT_FOUND", message, nil)
	}

	http.Redirect(w, r, urlRow.OriginalUrl, http.StatusFound)
}
