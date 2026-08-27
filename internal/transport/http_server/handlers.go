package httpserver

import (
	"encoding/json"
	"linkshortener/internal/apperrors"
	"linkshortener/internal/services"
	"log/slog"
	"net/http"
)

type Handler struct {
	logger  *slog.Logger
	service *services.LinkService
}

func NewHandler(logger *slog.Logger, service *services.LinkService) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

type CreateShortLinkRequest struct {
	OriginLink string `json:"origin_url"`
}

func (h *Handler) CreateShortLink(w http.ResponseWriter, r *http.Request) {
	var req CreateShortLinkRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Incorrect JSON format", http.StatusBadRequest)
		return
	}

	link, err := h.service.CreateLink(r.Context(), req.OriginLink)
	if err != nil {
		h.logger.Error("failed to create a link", "err", err, "url", req.OriginLink)
		http.Error(w, "Failed to creat a link", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(link); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}

}

func (h *Handler) RedirectToOriginal(w http.ResponseWriter, r *http.Request) {
	shortenCode := r.PathValue("code")

	if shortenCode == "" {
		http.Error(w, "URL is empty", http.StatusBadRequest)
	}

	original, err := h.service.GetByShorten(r.Context(), shortenCode)
	if err != nil {
		if err == apperrors.ErrLinkNotFound {
			http.Error(w, "Couldn't find original link", http.StatusNotFound)
		}
		h.logger.Error("Couldn't find original link", "error", err, "url", original)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, original, http.StatusTemporaryRedirect)
}
