package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/franpfeiffer/feriados-scraper/internal/cache"
	"github.com/franpfeiffer/feriados-scraper/internal/scraper"
)

var (
	feriadoScraper *scraper.FeriadoScraper
	feriadoCache   *cache.Cache
)

func main() {
	feriadoScraper = scraper.NewFeriadoScraper()
	feriadoCache = cache.NewCache(24 * time.Hour)

	http.HandleFunc("/api/feriados", handleFeriados)
	http.HandleFunc("/api/health", handleHealth)

	port := ":8080"
	log.Printf("server: http://localhost%s", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleFeriados(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "not allowed", http.StatusMethodNotAllowed)
		return
	}

	feriados, found := feriadoCache.Get()
	if !found {
		log.Println("no cache, scraping...")
		var err error
		feriados, err = feriadoScraper.GetFeriados()
		if err != nil {
			log.Printf("error grabbing feriados: %v", err)
			http.Error(w, fmt.Sprintf("error grabbing feriados: %v", err), http.StatusInternalServerError)
			return
		}

		feriadoCache.Set(feriados)
		log.Printf("feriados grabbed and stored in cache: %d feriados", len(feriados))
	} else {
		log.Println("cache hit")
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(map[string]any{
		"feriados": feriados,
		"total":    len(feriados),
	})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

