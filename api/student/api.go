package student

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"

	"github.com/dhax/go-base/auth/jwt"
	"github.com/dhax/go-base/logging"
	"github.com/dhax/go-base/models"
)

// Resource defines the student management resource
type Resource struct {
	Store     StudentStore
	AuthStore AuthTokenStore
}

// StudentStore defines database operations for student management
type StudentStore interface {
	GetStudentByID(ctx context.Context, id int64) (*models.Student, error)
	GetStudentByCustomUserID(ctx context.Context, customUserID int64) (*models.Student, error)
	CreateStudent(ctx context.Context, student *models.Student) error
	UpdateStudent(ctx context.Context, student *models.Student) error
	DeleteStudent(ctx context.Context, id int64) error
	ListStudents(ctx context.Context, filters map[string]interface{}) ([]models.Student, error)
	UpdateStudentLocation(ctx context.Context, id int64, locations map[string]bool) error
	CreateStudentVisit(ctx context.Context, studentID, roomID, timespanID int64) (*models.Visit, error)
	GetStudentVisits(ctx context.Context, studentID int64, date *time.Time) ([]models.Visit, error)
	GetRoomVisits(ctx context.Context, roomID int64, date *time.Time, active bool) ([]models.Visit, error)
	GetStudentAsList(ctx context.Context, id int64) (*models.StudentList, error)
	CreateFeedback(ctx context.Context, studentID int64, feedbackValue string, mensaFeedback bool) (*models.Feedback, error)
}

// AuthTokenStore defines operations for the auth token store
type AuthTokenStore interface {
	GetToken(t string) (*jwt.Token, error)
}

// NewResource creates a new student management resource
func NewResource(store StudentStore, authStore AuthTokenStore) *Resource {
	return &Resource{
		Store:     store,
		AuthStore: authStore,
	}
}

// Router creates a router for student management
func (rs *Resource) Router() chi.Router {
	r := chi.NewRouter()

	// JWT protected routes
	r.Group(func(r chi.Router) {
		r.Use(jwt.Authenticator)

		// Student routes
		r.Route("/", func(r chi.Router) {
			r.Get("/", rs.listStudents)
			r.Post("/", rs.createStudent)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", rs.getStudent)
				r.Put("/", rs.updateStudent)
				r.Delete("/", rs.deleteStudent)
				r.Get("/visits", rs.getStudentVisits)
			})
		})

		// Special operations
		r.Post("/register-in-room", rs.registerStudentInRoom)
		r.Post("/unregister-from-room", rs.unregisterStudentFromRoom)
		r.Post("/update-location", rs.updateStudentLocation)
		r.Post("/give-feedback", rs.giveFeedback)
	})

	return r
}

// StudentRequest is the request payload for Student data
type StudentRequest struct {
	*models.Student
}

// Bind preprocesses a StudentRequest
func (req *StudentRequest) Bind(r *http.Request) error {
	return nil
}

// LocationUpdateRequest is the request payload for updating a student's location
type LocationUpdateRequest struct {
	StudentID  int64 `json:"student_id"`
	InHouse    *bool `json:"in_house,omitempty"`
	WC         *bool `json:"wc,omitempty"`
	SchoolYard *bool `json:"school_yard,omitempty"`
}

// Bind preprocesses a LocationUpdateRequest
func (req *LocationUpdateRequest) Bind(r *http.Request) error {
	return nil
}

// RoomRegistrationRequest is the request payload for registering a student in a room
type RoomRegistrationRequest struct {
	StudentID int64  `json:"student_id"`
	DeviceID  string `json:"device_id"`
}

// Bind preprocesses a RoomRegistrationRequest
func (req *RoomRegistrationRequest) Bind(r *http.Request) error {
	return nil
}

// FeedbackRequest is the request payload for student feedback
type FeedbackRequest struct {
	StudentID     int64  `json:"student_id"`
	FeedbackValue string `json:"feedback_value"`
	MensaFeedback bool   `json:"mensa_feedback"`
}

// Bind preprocesses a FeedbackRequest
func (req *FeedbackRequest) Bind(r *http.Request) error {
	return nil
}

// VisitFilter is used to filter visit records
type VisitFilter struct {
	Date   *time.Time `json:"date,omitempty"`
	Active bool       `json:"active,omitempty"`
}

// Bind preprocesses a VisitFilter
func (req *VisitFilter) Bind(r *http.Request) error {
	return nil
}

// ======== Student Handlers ========

// listStudents returns a list of all students with optional filtering
func (rs *Resource) listStudents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters for filtering
	filters := make(map[string]interface{})

	if groupIDStr := r.URL.Query().Get("group_id"); groupIDStr != "" {
		groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
		if err == nil {
			filters["group_id"] = groupID
		}
	}

	if searchTerm := r.URL.Query().Get("search"); searchTerm != "" {
		filters["search"] = searchTerm
	}

	if inHouseStr := r.URL.Query().Get("in_house"); inHouseStr != "" {
		inHouse := inHouseStr == "true"
		filters["in_house"] = inHouse
	}

	students, err := rs.Store.ListStudents(ctx, filters)
	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	render.JSON(w, r, students)
}

// createStudent creates a new student
func (rs *Resource) createStudent(w http.ResponseWriter, r *http.Request) {
	data := &StudentRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	ctx := r.Context()
	if err := rs.Store.CreateStudent(ctx, data.Student); err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, data.Student)
}

// getStudent returns a specific student
func (rs *Resource) getStudent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(errors.New("invalid ID format")))
		return
	}

	ctx := r.Context()
	student, err := rs.Store.GetStudentByID(ctx, id)
	if err != nil {
		render.Render(w, r, ErrNotFound)
		return
	}

	render.JSON(w, r, student)
}

// updateStudent updates a specific student
func (rs *Resource) updateStudent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(errors.New("invalid ID format")))
		return
	}

	data := &StudentRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	ctx := r.Context()
	student, err := rs.Store.GetStudentByID(ctx, id)
	if err != nil {
		render.Render(w, r, ErrNotFound)
		return
	}

	// Update student fields, preserving what shouldn't be changed
	student.SchoolClass = data.SchoolClass
	student.Bus = data.Bus
	student.NameLG = data.NameLG
	student.ContactLG = data.ContactLG
	student.GroupID = data.GroupID

	if err := rs.Store.UpdateStudent(ctx, student); err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	render.JSON(w, r, student)
}

// deleteStudent deletes a specific student
func (rs *Resource) deleteStudent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(errors.New("invalid ID format")))
		return
	}

	ctx := r.Context()
	if err := rs.Store.DeleteStudent(ctx, id); err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	render.NoContent(w, r)
}

// ======== Special Operations ========

// updateStudentLocation updates a student's location flags
func (rs *Resource) updateStudentLocation(w http.ResponseWriter, r *http.Request) {
	data := &LocationUpdateRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	ctx := r.Context()

	// Create a map of location flags to update
	locations := make(map[string]bool)
	if data.InHouse != nil {
		locations["in_house"] = *data.InHouse
	}
	if data.WC != nil {
		locations["wc"] = *data.WC
	}
	if data.SchoolYard != nil {
		locations["school_yard"] = *data.SchoolYard
	}

	if err := rs.Store.UpdateStudentLocation(ctx, data.StudentID, locations); err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	student, err := rs.Store.GetStudentByID(ctx, data.StudentID)
	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	render.JSON(w, r, student)
}

// registerStudentInRoom registers a student in a room
func (rs *Resource) registerStudentInRoom(w http.ResponseWriter, r *http.Request) {
	data := &RoomRegistrationRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	ctx := r.Context()

	// Get the active room occupancy for the device ID
	log := logging.GetLogEntry(r)
	log.WithFields(logrus.Fields{
		"student_id": data.StudentID,
		"device_id":  data.DeviceID,
	}).Info("Registering student in room")

	// TODO: Get room occupancy for device ID
	// For now, we'll hardcode some values for demo purposes
	roomID := int64(1)
	timespanID := int64(1)

	// Create a visit record
	visit, err := rs.Store.CreateStudentVisit(ctx, data.StudentID, roomID, timespanID)
	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	// Update in_house status
	locations := map[string]bool{
		"in_house": true,
	}
	if err := rs.Store.UpdateStudentLocation(ctx, data.StudentID, locations); err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	render.JSON(w, r, visit)
}

// unregisterStudentFromRoom unregisters a student from a room
func (rs *Resource) unregisterStudentFromRoom(w http.ResponseWriter, r *http.Request) {
	data := &RoomRegistrationRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	ctx := r.Context()

	// Update in_house status to false
	locations := map[string]bool{
		"in_house": false,
	}
	if err := rs.Store.UpdateStudentLocation(ctx, data.StudentID, locations); err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	// Return success message
	render.JSON(w, r, map[string]interface{}{
		"success": true,
		"message": "Student unregistered from room",
	})
}

// getStudentVisits returns visits for a student
func (rs *Resource) getStudentVisits(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(errors.New("invalid ID format")))
		return
	}

	ctx := r.Context()

	// Parse date parameter
	var date *time.Time
	if dateStr := r.URL.Query().Get("date"); dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			date = &parsedDate
		}
	}

	visits, err := rs.Store.GetStudentVisits(ctx, id, date)
	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	render.JSON(w, r, visits)
}

// giveFeedback records feedback from a student
func (rs *Resource) giveFeedback(w http.ResponseWriter, r *http.Request) {
	data := &FeedbackRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	ctx := r.Context()
	feedback, err := rs.Store.CreateFeedback(ctx, data.StudentID, data.FeedbackValue, data.MensaFeedback)
	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, feedback)
}
