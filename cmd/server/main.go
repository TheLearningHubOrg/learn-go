package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/TheLearningHubOrg/learn-go/pkg/concurrency"
	"github.com/TheLearningHubOrg/learn-go/pkg/models"
)

var (
	itemsLock sync.RWMutex
	items     = []models.Item{
		{ID: 1, Title: "GMP Concurrency Engine", Payload: "Goroutine Work-Stealing", CreatedAt: time.Now()},
		{ID: 2, Title: "Low-Latency Tri-Color GC", Payload: "Sub-millisecond pauses", CreatedAt: time.Now()},
	}
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status := models.HealthStatus{
		Status:     "healthy",
		Version:    "1.23.0",
		Goroutines: runtime.NumGoroutine(),
		Timestamp:  time.Now(),
	}
	json.NewEncoder(w).Encode(status)
}

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		itemsLock.RLock()
		defer itemsLock.RUnlock()
		json.NewEncoder(w).Encode(items)
	case http.MethodPost:
		var newItem models.Item
		if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		itemsLock.Lock()
		newItem.ID = len(items) + 1
		newItem.CreatedAt = time.Now()
		items = append(items, newItem)
		itemsLock.Unlock()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newItem)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/api/v1/items", itemsHandler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Initialize background worker pool
	pool := concurrency.NewPool(4, 50)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pool.Start(ctx)

	// Graceful shutdown handling
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("⚡ Go Systems Masterclass Server listening on http://localhost:8080 (Goroutines: %d)", runtime.NumGoroutine())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server listen failed: %v", err)
		}
	}()

	<-stopChan
	log.Println("🛑 Shutting down server gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	pool.Stop()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server terminated cleanly.")
}
