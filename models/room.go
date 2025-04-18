package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"

	"github.com/uptrace/bun"
)

// Room represents a physical location in the facility
type Room struct {
	bun.BaseModel `bun:"table:rooms,alias:r"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	RoomName  string    `bun:"room_name,notnull,unique" json:"room_name"`
	Building  string    `bun:"building" json:"building,omitempty"`
	Floor     int       `bun:"floor,notnull,default:0" json:"floor"`
	Capacity  int       `bun:"capacity,notnull" json:"capacity"`
	Category  string    `bun:"category,notnull,default:'Other'" json:"category"`
	Color     string    `bun:"color,notnull,default:'#FFFFFF'" json:"color"`
	CreatedAt time.Time `bun:"created_at,notnull" json:"created_at,omitempty"`
	UpdatedAt time.Time `bun:"updated_at,notnull" json:"updated_at,omitempty"`
}

// BeforeInsert hook executed before database insert operation.
func (r *Room) BeforeInsert(db *bun.DB) error {
	now := time.Now()
	r.CreatedAt = now
	r.UpdatedAt = now
	return r.Validate()
}

// BeforeUpdate hook executed before database update operation.
func (r *Room) BeforeUpdate(db *bun.DB) error {
	r.UpdatedAt = time.Now()
	return r.Validate()
}

// Validate validates Room struct and returns validation errors.
func (r *Room) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.RoomName, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.Building, validation.Length(0, 100)),
		validation.Field(&r.Floor, validation.Min(0)),
		validation.Field(&r.Capacity, validation.Required, validation.Min(1)),
		validation.Field(&r.Category, validation.Required, validation.Length(1, 50)),
		validation.Field(&r.Color, validation.Required, is.HexColor),
	)
}
