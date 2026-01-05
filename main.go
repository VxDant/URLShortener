package main

import (
	"URLShortener/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := runMigrations(ctx); err != nil {
		log.Fatal("Migration error:", err)
	}

	log.Println("Database ready, starting server on port 8080...")

	mux := http.NewServeMux()

	store := NewInMemoryStore()
	service := NewURLService(store)

	mux.HandleFunc("GET /", homePage)
	mux.HandleFunc("GET /api/v1/shorturl/urls", service.getAllURL)
	mux.HandleFunc("POST /api/v1/shorturl/url", service.addURL)
	mux.HandleFunc("GET /{id}", service.navigatetoUrl)

	fmt.Println("Server starting on port 8080...")

	log.Fatal(http.ListenAndServe(":8080", mux))

	fmt.Println("Welcome to the url shortener app")

}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Simple REST API in Go!")
}

func (s *URLService) getAllURL(writer http.ResponseWriter, request *http.Request) {

	fmt.Printf("{ {"+
		"GET request, getAllURL} %v\n", s.getAllURLs())

	jsonData, error := json.Marshal(s.getAllURLs())

	if error != nil {
		fmt.Println(error)
		return
	}

	writer.WriteHeader(http.StatusOK)
	err := json.NewEncoder(writer).Encode(jsonData)

	if err != nil {
		fmt.Println(err)
		return
	}
}

func (s *URLService) addURL(writer http.ResponseWriter, request *http.Request) {
	longUrl := model.LongURL{URL: request.FormValue("url")}
	shortUrl := s.ProcessAndAddLongURLtoMap(longUrl)

	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(shortUrl)
}

func (s *URLService) navigatetoUrl(writer http.ResponseWriter, request *http.Request) {
	id := request.PathValue("id")

	longUrl := s.store.GetLongURL(model.ShortURL{Url: id})

	longUrl.URL = "https://" + longUrl.URL

	fmt.Printf("short url id: %v, long url: %v\n", id, longUrl.URL)

	//http.Redirect(writer, request, longUrl.URL, http.StatusFound)

	writer.Header().Set("Location", longUrl.URL)
	writer.WriteHeader(http.StatusFound)

}
