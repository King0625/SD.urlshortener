package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/King0625/SD.urlshortener/internal/service"
	"github.com/King0625/SD.urlshortener/pkg/utils"
	"github.com/go-chi/chi/v5"
)

type UrlHandler struct {
	Service service.UrlService
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

	result, err := h.Service.ShortenUrl(r.Context(), req.URL)

	if err != nil {
		fmt.Println(err)
		message = "create short url failed"
		utils.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", message, nil)
		return
	}

	res := ShortenResponse{
		ShortURL: os.Getenv("DOMAIN") + "/" + result.Code,
	}

	message = "Create short url successfully"
	utils.RespondSuccess(w, http.StatusCreated, message, res)
}

func (h *UrlHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	var message string
	code := chi.URLParam(r, "code")

	urlRow, err := h.Service.GetUrlByCode(r.Context(), code)
	if err != nil {
		message = "short url not found"
		utils.RespondError(w, http.StatusNotFound, "SHORT_URL_NOT_FOUND", message, nil)
		return
	}

	http.Redirect(w, r, urlRow.OriginalUrl, http.StatusFound)
}

func (h *UrlHandler) DeleteUrlByCode(w http.ResponseWriter, r *http.Request) {
	var message string
	code := chi.URLParam(r, "code")

	err := h.Service.DeleteUrlByCode(r.Context(), code)
	if err != nil {
		message = "delete url error"
		utils.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", message, nil)
		return
	}

	message = "delete short url successfully"
	utils.RespondSuccess(w, http.StatusOK, message, nil)
}
