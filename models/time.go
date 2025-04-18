package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/uptrace/bun"
)

// Timespan represents a duration with start and optional end time
type Timespan struct {
	bun.BaseModel `bun:"table:timespans,alias:ts"`

	ID        int64      `bun:"id,pk,autoincrement" json:"id"`
	StartTime time.Time  `bun:"starttime,notnull" json:"starttime"`
	EndTime   *time.Time `bun:"endtime" json:"endtime,omitempty"`
	CreatedAt time.Time  `bun:"created_at,notnull" json:"created_at,omitempty"`
}

// BeforeInsert hook executed before database insert operation.
func (t *Timespan) BeforeInsert(db *bun.DB) error {
	t.CreatedAt = time.Now()
	return t.Validate()
}

// Validate validates Timespan struct and returns validation errors.
func (t *Timespan) Validate() error {
	return validation.ValidateStruct(t,
		validation.Field(&t.StartTime, validation.Required),
	)
}

// Datespan represents a period between two dates
type Datespan struct {
	bun.BaseModel `bun:"table:datespans,alias:ds"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	StartDate time.Time `bun:"startdate,notnull" json:"startdate"`
	EndDate   time.Time `bun:"enddate,notnull" json:"enddate"`
	CreatedAt time.Time `bun:"created_at,notnull" json:"created_at,omitempty"`
}

// BeforeInsert hook executed before database insert operation.
func (d *Datespan) BeforeInsert(db *bun.DB) error {
	d.CreatedAt = time.Now()
	return d.Validate()
}

// Validate validates Datespan struct and returns validation errors.
func (d *Datespan) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.StartDate, validation.Required),
		validation.Field(&d.EndDate, validation.Required),
	)
}
