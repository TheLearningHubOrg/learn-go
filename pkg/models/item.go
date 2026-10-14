package models

import "time"

// Item represents a high-throughput entity in the distributed store
type Item struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Payload   string    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

// HealthStatus represents the runtime health metrics
type HealthStatus struct {
	Status    string    `json:"status"`
	Version   string    `json:"version"`
	Goroutines int      `json:"goroutines"`
	Timestamp time.Time `json:"timestamp"`
}
