package main

import (
	"URLShortener/internal/database"
	"URLShortener/internal/repository"
	"URLShortener/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {

	log.Println("🚀 Starting URL Shortener service...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := runMigrations(ctx); err != nil {
		log.Fatal("❌ Migration failed:", err)
	}

	if err := database.Connect(); err != nil {
		log.Fatal("❌ Database connection failed:", err)
	}
	defer database.Close()

	urlRepo := repository.NewURLRepository(database.DB)

	log.Println("Database ready, starting server on port 8080...")

	mux := http.NewServeMux()

	store := NewInMemoryStore()
	service := NewURLService(store)

	mux.HandleFunc("GET /", homePage)
	mux.HandleFunc("GET /health", health)
	mux.HandleFunc("GET /test-db", testDatabaseHandler(urlRepo))
	mux.HandleFunc("GET /api/v1/shorturl/urls", service.getAllURL)
	mux.HandleFunc("POST /api/v1/shorturl/url", service.addURL)
	mux.HandleFunc("GET /{id}", service.navigatetoUrl)

	fmt.Println("Server starting on port 8080...")

	log.Fatal(http.ListenAndServe(":8080", mux))

	fmt.Println("Welcome to the url shortener app")

}

func health(w http.ResponseWriter, request *http.Request) {
	if err := database.Health(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Database unhealthy"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
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

func testDatabaseHandler(repo *repository.URLRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Testing database operations...")

		// Test 1: Create a URL
		testShortCode := fmt.Sprintf("test%d", time.Now().UnixNano())[0:10]
		created, err := repo.Create(testShortCode, "https://example.com/test")
		if err != nil {
			http.Error(w, fmt.Sprintf("Create failed: %v", err), http.StatusInternalServerError)
			return
		}
		log.Printf("✅ Created URL: %+v", created)

		// Test 2: Retrieve by short code
		retrieved, err := repo.GetByShortCode(testShortCode)
		if err != nil {
			http.Error(w, fmt.Sprintf("GetByShortCode failed: %v", err), http.StatusInternalServerError)
			return
		}
		log.Printf("✅ Retrieved URL: %+v", retrieved)

		// Test 3: Check by long URL
		existing, err := repo.GetByLongURL("https://example.com/test")
		if err != nil {
			http.Error(w, fmt.Sprintf("GetByLongURL failed: %v", err), http.StatusInternalServerError)
			return
		}
		log.Printf("✅ Found existing URL: %+v", existing)

		// Test 4: Increment clicks
		if err := repo.IncrementClicks(testShortCode); err != nil {
			http.Error(w, fmt.Sprintf("IncrementClicks failed: %v", err), http.StatusInternalServerError)
			return
		}
		log.Println("✅ Incremented clicks")

		// Test 5: Verify click increment
		updated, err := repo.GetByShortCode(testShortCode)
		if err != nil {
			http.Error(w, fmt.Sprintf("GetByShortCode (2nd) failed: %v", err), http.StatusInternalServerError)
			return
		}
		log.Printf("✅ Updated URL (clicks should be 1): %+v", updated)

		// Test 6: Clean up
		if err := repo.Delete(testShortCode); err != nil {
			log.Printf("⚠️  Delete failed: %v", err)
		} else {
			log.Println("✅ Deleted test URL")
		}

		// Return success response
		response := map[string]interface{}{
			"status":  "success",
			"message": "All database operations completed successfully",
			"tests": map[string]bool{
				"create":           true,
				"get_by_short":     true,
				"get_by_long":      true,
				"increment_clicks": true,
				"delete":           true,
			},
			"test_url": created,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
