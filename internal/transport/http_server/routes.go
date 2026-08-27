package httpserver

import "net/http"

func (h *Handler) Route(mux *http.ServeMux) {
	mux.HandleFunc("POST /shorten/links", h.CreateShortLink)
	mux.HandleFunc("GET /shorten/{code}", h.RedirectToOriginal)
}
