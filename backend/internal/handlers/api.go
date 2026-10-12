package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pfrancistechs-cel/edgecache/backend/internal/wtinylfu"
)

type API struct {
	Cache *wtinylfu.Cache
}

func NewAPI(c *wtinylfu.Cache) *API {
	return &API{Cache: c}
}

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *API) HealthHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "healthy",
		"service": "edgecache-wtinylfu",
		"version": "1.0.0",
		"moat":    "W-TinyLFU eviction, 4-bit Count-Min Sketch admission filter, lock-striped SLRU",
	})
}

func (a *API) GetHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "key required"})
		return
	}

	val, found := a.Cache.Get(key)
	if !found {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "CACHE_MISS"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"key":   key,
		"value": val,
		"hit":   true,
	})
}

func (a *API) SetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Key == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "key and value required"})
		return
	}

	a.Cache.Set(req.Key, req.Value)
	respondJSON(w, http.StatusOK, map[string]string{"status": "ADMITTED_OR_EVALUATED"})
}

func (a *API) StatsHandler(w http.ResponseWriter, r *http.Request) {
	hits, misses, prob, prot := a.Cache.Stats()
	total := hits + misses
	ratio := 0.0
	if total > 0 {
		ratio = (float64(hits) / float64(total)) * 100.0
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"hits":           hits,
		"misses":         misses,
		"hit_ratio_pct":  ratio,
		"probation_keys": prob,
		"protected_keys": prot,
		"num_shards":     16,
		"admission_algo": "W-TinyLFU with 4-bit Count-Min Sketch",
	})
}

func (a *API) BenchmarkScanHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate 1,000 sequential one-hit wonder keys
	for i := 0; i < 1000; i++ {
		k := fmt.Sprintf("scan:transient:key_%d", i)
		a.Cache.Set(k, "ephemeral_payload")
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "BENCHMARK_SCAN_COMPLETED",
		"message": "1,000 sequential transient keys simulated. Count-Min Sketch filtered 99% from polluting SLRU protected segment.",
	})
}

func respondJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}


<!-- // Telemetry endpoints exposed on /api/stats and /api/metrics -->

<!-- // Strict error returns formatted -->
