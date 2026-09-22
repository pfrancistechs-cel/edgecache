package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pfrancistechs-cel/edgecache/backend/internal/handlers"
	"github.com/pfrancistechs-cel/edgecache/backend/internal/wtinylfu"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8088"
	}

	c := wtinylfu.NewCache(4000, 16)

	// Seed hot records
	for i := 1; i <= 20; i++ {
		k := fmt.Sprintf("user:session:hot_%02d", i)
		c.Set(k, fmt.Sprintf(`{"role":"admin","uid":%d}`, 1000+i))
		for j := 0; j < 5; j++ {
			c.Get(k) // Build frequency in Count-Min Sketch & promote to Protected SLRU
		}
	}

	api := handlers.NewAPI(c)
	mux := http.NewServeMux()

	mux.HandleFunc("/health", api.HealthHandler)
	mux.HandleFunc("/api/v1/get", api.GetHandler)
	mux.HandleFunc("/api/v1/set", api.SetHandler)
	mux.HandleFunc("/api/v1/stats", api.StatsHandler)
	mux.HandleFunc("/api/v1/benchmark-scan", api.BenchmarkScanHandler)

	handler := handlers.EnableCORS(mux)

	serverAddr := ":" + port
	fmt.Printf("⚡ EdgeCache W-TinyLFU Engine running on http://localhost%s\n", serverAddr)
	if err := http.ListenAndServe(serverAddr, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

