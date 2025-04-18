package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/uptrace/bun"
)

// Student represents a student in the system
type Student struct {
	bun.BaseModel `bun:"table:students,alias:s"`

	ID           int64       `bun:"id,pk,autoincrement" json:"id"`
	SchoolClass  string      `bun:"school_class,notnull" json:"school_class"`
	Bus          bool        `bun:"bus,notnull,default:false" json:"bus"`
	NameLG       string      `bun:"name_lg,notnull" json:"name_lg"`
	ContactLG    string      `bun:"contact_lg,notnull" json:"contact_lg"`
	InHouse      bool        `bun:"in_house,notnull,default:false" json:"in_house"`
	WC           bool        `bun:"wc,notnull,default:false" json:"wc"`
	SchoolYard   bool        `bun:"school_yard,notnull,default:false" json:"school_yard"`
	CustomUserID int64       `bun:"custom_user_id,notnull" json:"custom_user_id"`
	GroupID      int64       `bun:"group_id,notnull" json:"group_id"`
	CreatedAt    time.Time   `bun:"created_at,notnull" json:"created_at,omitempty"`
	UpdatedAt    time.Time   `bun:"updated_at,notnull" json:"updated_at,omitempty"`
	CustomUser   *CustomUser `bun:"rel:belongs-to,join:custom_user_id=id" json:"custom_user,omitempty"`
	Group        *Group      `bun:"rel:belongs-to,join:group_id=id" json:"group,omitempty"`
}

// BeforeInsert hook executed before database insert operation.
func (s *Student) BeforeInsert(db *bun.DB) error {
	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now
	return s.Validate()
}

// BeforeUpdate hook executed before database update operation.
func (s *Student) BeforeUpdate(db *bun.DB) error {
	s.UpdatedAt = time.Now()
	return s.Validate()
}

// Validate validates Student struct and returns validation errors.
func (s *Student) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.SchoolClass, validation.Required, validation.Length(1, 50)),
		validation.Field(&s.NameLG, validation.Required, validation.Length(1, 100)),
		validation.Field(&s.ContactLG, validation.Required, validation.Length(1, 100)),
		validation.Field(&s.CustomUserID, validation.Required),
		validation.Field(&s.GroupID, validation.Required),
	)
}

// Visit tracks a student's presence in a room
type Visit struct {
	bun.BaseModel `bun:"table:visits,alias:v"`

	ID         int64     `bun:"id,pk,autoincrement" json:"id"`
	Day        time.Time `bun:"day,notnull" json:"day"`
	StudentID  int64     `bun:"student_id,notnull" json:"student_id"`
	RoomID     int64     `bun:"room_id,notnull" json:"room_id"`
	TimespanID int64     `bun:"timespan_id,notnull" json:"timespan_id"`
	CreatedAt  time.Time `bun:"created_at,notnull" json:"created_at,omitempty"`
	Student    *Student  `bun:"rel:belongs-to,join:student_id=id" json:"student,omitempty"`
	Room       *Room     `bun:"rel:belongs-to,join:room_id=id" json:"room,omitempty"`
}

// BeforeInsert hook executed before database insert operation.
func (v *Visit) BeforeInsert(db *bun.DB) error {
	v.CreatedAt = time.Now()
	return v.Validate()
}

// Validate validates Visit struct and returns validation errors.
func (v *Visit) Validate() error {
	return validation.ValidateStruct(v,
		validation.Field(&v.Day, validation.Required),
		validation.Field(&v.StudentID, validation.Required),
		validation.Field(&v.RoomID, validation.Required),
		validation.Field(&v.TimespanID, validation.Required),
	)
}

// Feedback stores student feedback
type Feedback struct {
	bun.BaseModel `bun:"table:feedback,alias:f"`

	ID            int64     `bun:"id,pk,autoincrement" json:"id"`
	FeedbackValue string    `bun:"feedback_value,notnull" json:"feedback_value"`
	Day           time.Time `bun:"day,notnull" json:"day"`
	Time          time.Time `bun:"time,notnull" json:"time"`
	StudentID     int64     `bun:"student_id,notnull" json:"student_id"`
	MensaFeedback bool      `bun:"mensa_feedback,notnull,default:false" json:"mensa_feedback"`
	CreatedAt     time.Time `bun:"created_at,notnull" json:"created_at,omitempty"`
	Student       *Student  `bun:"rel:belongs-to,join:student_id=id" json:"student,omitempty"`
}

// BeforeInsert hook executed before database insert operation.
func (f *Feedback) BeforeInsert(db *bun.DB) error {
	f.CreatedAt = time.Now()
	return f.Validate()
}

// Validate validates Feedback struct and returns validation errors.
func (f *Feedback) Validate() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.FeedbackValue, validation.Required),
		validation.Field(&f.Day, validation.Required),
		validation.Field(&f.Time, validation.Required),
		validation.Field(&f.StudentID, validation.Required),
	)
}
