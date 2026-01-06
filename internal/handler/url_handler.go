package handler

import (
	"URLShortener/internal/service"
	"encoding/json"
	"log"
	"net/http"
)

type URLHandler struct {
	service *service.URLService
}

func NewURLHandler(service *service.URLService) *URLHandler {
	return &URLHandler{service: service}
}

func (h *URLHandler) GetAllURLs(w http.ResponseWriter, r *http.Request) {
	urls, err := h.service.GetAllURL()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error fetching URLs: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(urls)
}

func (h *URLHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	t := struct {
		URL string `json:"url"`
	}{}

	if err := decoder.Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	url, err := h.service.CreateAndProcessShortURL(t.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error creating URL: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(url)
}
