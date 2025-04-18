package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/dhax/go-base/models"
	"github.com/uptrace/bun"
)

// StudentStore implements database operations for student management
type StudentStore struct {
	db *bun.DB
}

// NewStudentStore returns a StudentStore
func NewStudentStore(db *bun.DB) *StudentStore {
	return &StudentStore{
		db: db,
	}
}

// GetStudentByID retrieves a Student by ID with related CustomUser
func (s *StudentStore) GetStudentByID(ctx context.Context, id int64) (*models.Student, error) {
	student := new(models.Student)
	err := s.db.NewSelect().
		Model(student).
		Relation("CustomUser").
		Relation("Group").
		Where("student.id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return student, nil
}

// GetStudentByCustomUserID retrieves a Student by CustomUserID
func (s *StudentStore) GetStudentByCustomUserID(ctx context.Context, customUserID int64) (*models.Student, error) {
	student := new(models.Student)
	err := s.db.NewSelect().
		Model(student).
		Relation("CustomUser").
		Relation("Group").
		Where("student.custom_user_id = ?", customUserID).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return student, nil
}

// CreateStudent creates a new Student
func (s *StudentStore) CreateStudent(ctx context.Context, student *models.Student) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if a student already exists for this custom user
	exists, err := tx.NewSelect().
		Model((*models.Student)(nil)).
		Where("custom_user_id = ?", student.CustomUserID).
		Exists(ctx)

	if err != nil {
		return err
	}

	if exists {
		return errors.New("a student already exists for this custom user")
	}

	// Create the student
	_, err = tx.NewInsert().
		Model(student).
		Exec(ctx)

	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// UpdateStudent updates an existing Student
func (s *StudentStore) UpdateStudent(ctx context.Context, student *models.Student) error {
	_, err := s.db.NewUpdate().
		Model(student).
		WherePK().
		Exec(ctx)

	return err
}

// DeleteStudent deletes a Student
func (s *StudentStore) DeleteStudent(ctx context.Context, id int64) error {
	_, err := s.db.NewDelete().
		Model((*models.Student)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	return err
}

// ListStudents returns a list of all Students with related CustomUser and Group
func (s *StudentStore) ListStudents(ctx context.Context, filters map[string]interface{}) ([]models.Student, error) {
	var students []models.Student

	query := s.db.NewSelect().
		Model(&students).
		Relation("CustomUser").
		Relation("Group")

	// Apply filters
	if groupID, ok := filters["group_id"].(int64); ok && groupID > 0 {
		query = query.Where("student.group_id = ?", groupID)
	}

	if searchTerm, ok := filters["search"].(string); ok && searchTerm != "" {
		query = query.Where("custom_user.first_name ILIKE ? OR custom_user.second_name ILIKE ?",
			"%"+searchTerm+"%", "%"+searchTerm+"%")
	}

	if inHouse, ok := filters["in_house"].(bool); ok {
		query = query.Where("student.in_house = ?", inHouse)
	}

	err := query.OrderExpr("custom_user.first_name ASC").
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return students, nil
}

// UpdateStudentLocation updates a student's location flags (in_house, wc, school_yard)
func (s *StudentStore) UpdateStudentLocation(ctx context.Context, id int64, locations map[string]bool) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	student, err := s.GetStudentByID(ctx, id)
	if err != nil {
		return err
	}

	// Update student location flags
	if inHouse, ok := locations["in_house"]; ok {
		student.InHouse = inHouse
	}
	if wc, ok := locations["wc"]; ok {
		student.WC = wc
	}
	if schoolYard, ok := locations["school_yard"]; ok {
		student.SchoolYard = schoolYard
	}

	_, err = tx.NewUpdate().
		Model(student).
		Column("in_house", "wc", "school_yard", "modified_at").
		WherePK().
		Exec(ctx)

	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// CreateStudentVisit creates a new Visit record for a student
func (s *StudentStore) CreateStudentVisit(ctx context.Context, studentID, roomID, timespanID int64) (*models.Visit, error) {
	visit := &models.Visit{
		Day:        time.Now(),
		StudentID:  studentID,
		RoomID:     roomID,
		TimespanID: timespanID,
		CreatedAt:  time.Now(),
	}

	_, err := s.db.NewInsert().
		Model(visit).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	return visit, nil
}

// GetStudentVisits retrieves all visits for a student
func (s *StudentStore) GetStudentVisits(ctx context.Context, studentID int64, date *time.Time) ([]models.Visit, error) {
	var visits []models.Visit

	query := s.db.NewSelect().
		Model(&visits).
		Relation("Room").
		Relation("Timespan").
		Where("student_id = ?", studentID)

	if date != nil {
		query = query.Where("DATE(day) = DATE(?)", *date)
	}

	err := query.Order("day DESC").
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return visits, nil
}

// GetRoomVisits retrieves all visits for a room
func (s *StudentStore) GetRoomVisits(ctx context.Context, roomID int64, date *time.Time, active bool) ([]models.Visit, error) {
	var visits []models.Visit

	query := s.db.NewSelect().
		Model(&visits).
		Relation("Student").
		Relation("Student.CustomUser").
		Relation("Timespan").
		Where("room_id = ?", roomID)

	if date != nil {
		query = query.Where("DATE(day) = DATE(?)", *date)
	}

	if active {
		query = query.Where("timespan.endtime IS NULL")
	}

	err := query.Order("day DESC").
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return visits, nil
}

// GetStudentAsList returns a student as a StudentList object
func (s *StudentStore) GetStudentAsList(ctx context.Context, id int64) (*models.StudentList, error) {
	student, err := s.GetStudentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if student.CustomUser == nil || student.Group == nil {
		return nil, errors.New("student data incomplete")
	}

	studentList := &models.StudentList{
		ID:          student.ID,
		Name:        student.CustomUser.FirstName + " " + student.CustomUser.SecondName,
		SchoolClass: student.SchoolClass,
		GroupName:   student.Group.Name,
		InHouse:     student.InHouse,
	}

	return studentList, nil
}

// CreateFeedback creates a new Feedback record for a student
func (s *StudentStore) CreateFeedback(ctx context.Context, studentID int64, feedbackValue string, mensaFeedback bool) (*models.Feedback, error) {
	now := time.Now()
	feedback := &models.Feedback{
		FeedbackValue: feedbackValue,
		Day:           now,
		Time:          now,
		StudentID:     studentID,
		MensaFeedback: mensaFeedback,
		CreatedAt:     now,
	}

	_, err := s.db.NewInsert().
		Model(feedback).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	return feedback, nil
}
