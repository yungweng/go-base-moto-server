// Package rfid handles RFID tag communication from Raspberry Pi
package rfid

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/uptrace/bun"
)

// API provides RFID handlers.
type API struct {
	store RFIDStore
}

// NewAPI configures and returns RFID API.
func NewAPI(db *bun.DB) (*API, error) {
	store := NewRFIDStore(db)
	api := &API{
		store: store,
	}
	return api, nil
}

// Router provides RFID routes.
func (a *API) Router() *chi.Mux {
	r := chi.NewRouter()

	// Endpoints for RFID Python Daemon
	r.Post("/tag", a.handleTagRead)
	r.Get("/tags", a.handleGetAllTags)

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