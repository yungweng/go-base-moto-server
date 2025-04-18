package rfid

import (
	"net/http"
	"time"
)

// Tag represents an RFID tag read
type Tag struct {
	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	TagID     string    `json:"tag_id" bun:"tag_id,notnull"`
	ReaderID  string    `json:"reader_id" bun:"reader_id,notnull"`
	ReadAt    time.Time `json:"read_at" bun:"read_at,notnull"`
	CreatedAt time.Time `json:"created_at" bun:"created_at,notnull"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at,notnull"`
}

// TagReadRequest is the payload for tag read endpoint
type TagReadRequest struct {
	TagID    string `json:"tag_id"`
	ReaderID string `json:"reader_id"`
	ReadAt   string `json:"read_at,omitempty"` // Optional ISO8601 timestamp
}

// Bind preprocesses a TagReadRequest
func (t *TagReadRequest) Bind(r *http.Request) error {
	// Validation logic can be added here
	return nil
}

// TauriSyncRequest is the payload for Tauri app sync endpoint
type TauriSyncRequest struct {
	DeviceID string    `json:"device_id"`
	Data     []SyncTag `json:"data"`
}

// SyncTag represents a tag record from the Tauri app
type SyncTag struct {
	TagID      string    `json:"tag_id"`
	ReaderID   string    `json:"reader_id"`
	LocalReadAt time.Time `json:"local_read_at"`
}

// Bind preprocesses a TauriSyncRequest
func (t *TauriSyncRequest) Bind(r *http.Request) error {
	// Validation logic can be added here
	return nil
}

// TauriSyncResponse is the response for the Tauri app sync endpoint
type TauriSyncResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// AppStats contains statistics for the Tauri app
type AppStats struct {
	TagCount int `json:"tag_count"`
}

// AppStatus represents the server status for the Tauri app
type AppStatus struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Stats     AppStats  `json:"stats"`
}