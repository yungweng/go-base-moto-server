// Package rfid handles RFID tag communication from Raspberry Pi
package rfid

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"

	"github.com/dhax/go-base/logging"
	"github.com/dhax/go-base/models"
)

// API provides RFID handlers.
type API struct {
	store        RFIDStore
	userStore    UserStore
	studentStore StudentStore
}

// UserStore defines operations needed from the user store
type UserStore interface {
	GetCustomUserByTagID(ctx context.Context, tagID string) (*models.CustomUser, error)
}

// StudentStore defines operations needed from the student store
type StudentStore interface {
	GetStudentByCustomUserID(ctx context.Context, customUserID int64) (*models.Student, error)
	UpdateStudentLocation(ctx context.Context, id int64, locations map[string]bool) error
}

// NewAPI configures and returns RFID API.
func NewAPI(db *bun.DB) (*API, error) {
	store := NewRFIDStore(db)
	api := &API{
		store: store,
	}
	return api, nil
}

// SetUserStore sets the user store for RFID API
func (a *API) SetUserStore(userStore UserStore) {
	a.userStore = userStore
}

// SetStudentStore sets the student store for RFID API
func (a *API) SetStudentStore(studentStore StudentStore) {
	a.studentStore = studentStore
}

// Router provides RFID routes.
func (a *API) Router() *chi.Mux {
	r := chi.NewRouter()

	// Endpoints for RFID Python Daemon
	r.Post("/tag", a.handleTagRead)
	r.Get("/tags", a.handleGetAllTags)

	// Student tracking with RFID
	r.Post("/track-student", a.handleStudentTracking)

	// Endpoints for Tauri App
	r.Post("/app/sync", a.handleTauriSync)
	r.Get("/app/status", a.handleTauriStatus)

	return r
}

// handleTagRead processes RFID tag reads from the Raspberry Pi
func (a *API) handleTagRead(w http.ResponseWriter, r *http.Request) {
	data := &TagReadRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	tag, err := a.store.SaveTag(r.Context(), data.TagID, data.ReaderID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, tag)
}

// handleGetAllTags returns all stored RFID tags
func (a *API) handleGetAllTags(w http.ResponseWriter, r *http.Request) {
	tags, err := a.store.GetAllTags(r.Context())
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.JSON(w, r, tags)
}

// handleTauriSync processes synchronization requests from the Tauri app
func (a *API) handleTauriSync(w http.ResponseWriter, r *http.Request) {
	data := &TauriSyncRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	err := a.store.SaveTauriTags(r.Context(), data.DeviceID, data.Data)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	response := &TauriSyncResponse{
		Success: true,
		Message: fmt.Sprintf("Successfully synced %d tags", len(data.Data)),
	}

	render.JSON(w, r, response)
}

// handleTauriStatus returns status information to the Tauri app
func (a *API) handleTauriStatus(w http.ResponseWriter, r *http.Request) {
	tagCount, err := a.store.GetTagStats(r.Context())
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	status := &AppStatus{
		Status:    "ok",
		Timestamp: time.Now(),
		Stats: AppStats{
			TagCount: tagCount,
		},
	}

	render.JSON(w, r, status)
}

// StudentTrackingRequest is the request for student tracking with RFID
type StudentTrackingRequest struct {
	TagID        string `json:"tag_id"`
	ReaderID     string `json:"reader_id"`
	LocationType string `json:"location_type"` // "entry", "wc", "schoolyard", or "exit"
}

// Bind preprocesses a StudentTrackingRequest
func (req *StudentTrackingRequest) Bind(r *http.Request) error {
	return nil
}

// StudentTrackingResponse is the response for student tracking
type StudentTrackingResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	StudentID int64  `json:"student_id,omitempty"`
	Name      string `json:"name,omitempty"`
	Location  string `json:"location,omitempty"`
}

// handleStudentTracking processes RFID tags for student tracking
func (a *API) handleStudentTracking(w http.ResponseWriter, r *http.Request) {
	data := &StudentTrackingRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	ctx := r.Context()
	// Get the logger - our improved GetLogEntry will handle all the edge cases
	logger := logging.GetLogEntry(r)

	// Create an entry from the logger - it's always safe to do this
	log := logrus.NewEntry(logrus.StandardLogger())
	if fieldLogger, ok := logger.(*logrus.Entry); ok {
		log = fieldLogger
	} else if logger != nil {
		// If it's not nil but not a *logrus.Entry, create a new entry that copies the fields
		log = logrus.NewEntry(logrus.StandardLogger()).WithFields(logrus.Fields{
			"request_id": middleware.GetReqID(ctx),
			"method":     r.Method,
			"path":       r.URL.Path,
		})
	}

	// Check if userStore is set
	if a.userStore == nil {
		render.Render(w, r, ErrInternalServer(fmt.Errorf("user store not configured")))
		return
	}

	// First log the tag read
	_, err := a.store.SaveTag(ctx, data.TagID, data.ReaderID)
	if err != nil {
		log.WithError(err).Error("Failed to save tag read")
		// Continue anyway, as this is just for logging
	}

	// Find the user by tag ID
	user, err := a.userStore.GetCustomUserByTagID(ctx, data.TagID)
	if err != nil {
		log.WithFields(logrus.Fields{
			"tag_id":    data.TagID,
			"reader_id": data.ReaderID,
			"error":     err.Error(),
		}).Info("Tag scanned but no user found")

		render.JSON(w, r, &StudentTrackingResponse{
			Success: false,
			Message: "No user found with this tag ID",
		})
		return
	}

	// Determine the location update
	locationUpdate := ""
	locationUpdates := make(map[string]bool)

	switch data.LocationType {
	case "entry":
		locationUpdate = "in-house"
		locationUpdates["in_house"] = true
		locationUpdates["wc"] = false
		locationUpdates["school_yard"] = false
	case "wc":
		locationUpdate = "bathroom"
		locationUpdates["in_house"] = true
		locationUpdates["wc"] = true
		locationUpdates["school_yard"] = false
	case "schoolyard":
		locationUpdate = "schoolyard"
		locationUpdates["in_house"] = false
		locationUpdates["wc"] = false
		locationUpdates["school_yard"] = true
	case "exit":
		locationUpdate = "out"
		locationUpdates["in_house"] = false
		locationUpdates["wc"] = false
		locationUpdates["school_yard"] = false
	default:
		locationUpdate = "unknown"
	}

	// Update student location if studentStore is configured
	studentID := int64(0)
	if a.studentStore != nil {
		student, err := a.studentStore.GetStudentByCustomUserID(ctx, user.ID)
		if err == nil {
			studentID = student.ID
			err = a.studentStore.UpdateStudentLocation(ctx, student.ID, locationUpdates)
			if err != nil {
				log.WithError(err).Error("Failed to update student location")
			}
		} else {
			log.WithError(err).WithField("user_id", user.ID).Warning("Found user but not student record")
		}
	}

	// Log the tracking event
	log.WithFields(logrus.Fields{
		"tag_id":     data.TagID,
		"reader_id":  data.ReaderID,
		"user_id":    user.ID,
		"user_name":  user.FirstName + " " + user.SecondName,
		"location":   data.LocationType,
		"student_id": studentID,
	}).Info("Student location tracked")

	// Return the tracking response
	render.JSON(w, r, &StudentTrackingResponse{
		Success:   true,
		Message:   "Location tracking recorded",
		StudentID: user.ID,
		Name:      user.FirstName + " " + user.SecondName,
		Location:  locationUpdate,
	})
}
