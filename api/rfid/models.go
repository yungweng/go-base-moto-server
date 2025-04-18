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
	TagID       string    `json:"tag_id"`
	ReaderID    string    `json:"reader_id"`
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
	TagCount             int `json:"tag_count"`
	StudentsInHouse      int `json:"students_in_house"`
	StudentsInWC         int `json:"students_in_wc"`
	StudentsInSchoolYard int `json:"students_in_school_yard"`
}

// AppStatus represents the server status for the Tauri app
type AppStatus struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Stats     AppStats  `json:"stats"`
	Version   string    `json:"version"`
}

// RoomEntryRequest represents a student entering a room with an RFID tag
type RoomEntryRequest struct {
	TagID    string `json:"tag_id"`
	RoomID   int64  `json:"room_id"`
	ReaderID string `json:"reader_id"`
}

// Bind preprocesses a RoomEntryRequest
func (req *RoomEntryRequest) Bind(r *http.Request) error {
	return nil
}

// RoomExitRequest represents a student exiting a room with an RFID tag
type RoomExitRequest struct {
	TagID    string `json:"tag_id"`
	RoomID   int64  `json:"room_id"`
	ReaderID string `json:"reader_id"`
}

// Bind preprocesses a RoomExitRequest
func (req *RoomExitRequest) Bind(r *http.Request) error {
	return nil
}

// OccupancyResponse represents a response for room occupancy operations
type OccupancyResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	StudentID    int64  `json:"student_id,omitempty"`
	RoomID       int64  `json:"room_id,omitempty"`
	StudentCount int    `json:"student_count,omitempty"`
}

// RoomOccupancyStudent represents a student in a room for occupancy reporting
type RoomOccupancyStudent struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	EnteredAt time.Time `json:"entered_at"`
}

// RoomOccupancyData represents the occupancy data for a room
type RoomOccupancyData struct {
	RoomID       int64                  `json:"room_id"`
	RoomName     string                 `json:"room_name"`
	Capacity     int                    `json:"capacity"`
	StudentCount int                    `json:"student_count"`
	Students     []RoomOccupancyStudent `json:"students"`
}
