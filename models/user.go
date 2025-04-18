package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/uptrace/bun"
)

// CustomUser extends the base Account model with additional user information
type CustomUser struct {
	bun.BaseModel `bun:"table:custom_users,alias:cu"`

	ID         int64     `bun:"id,pk,autoincrement" json:"id"`
	FirstName  string    `bun:"first_name,notnull" json:"first_name"`
	SecondName string    `bun:"second_name,notnull" json:"second_name"`
	TagID      string    `bun:"tag_id,unique" json:"tag_id,omitempty"`
	AccountID  int64     `bun:"account_id,notnull" json:"account_id"`
	CreatedAt  time.Time `bun:"created_at,notnull" json:"created_at,omitempty"`
	UpdatedAt  time.Time `bun:"updated_at,notnull" json:"updated_at,omitempty"`
}

// BeforeInsert hook executed before database insert operation.
func (u *CustomUser) BeforeInsert(db *bun.DB) error {
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	return u.Validate()
}

// BeforeUpdate hook executed before database update operation.
func (u *CustomUser) BeforeUpdate(db *bun.DB) error {
	u.UpdatedAt = time.Now()
	return u.Validate()
}

// Validate validates CustomUser struct and returns validation errors.
func (u *CustomUser) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.FirstName, validation.Required, validation.Length(1, 100)),
		validation.Field(&u.SecondName, validation.Required, validation.Length(1, 100)),
		validation.Field(&u.TagID, validation.Length(0, 50)),
		validation.Field(&u.AccountID, validation.Required),
	)
}

// PedagogicalSpecialist represents a staff member with specific educational roles
type PedagogicalSpecialist struct {
	bun.BaseModel `bun:"table:pedagogical_specialists,alias:ps"`

	ID            int64       `bun:"id,pk,autoincrement" json:"id"`
	Role          string      `bun:"role,notnull" json:"role"`
	CustomUserID  int64       `bun:"custom_user_id,notnull" json:"custom_user_id"`
	UserID        int64       `bun:"user_id,notnull,unique" json:"user_id"`
	IsPasswordOTP bool        `bun:"is_password_otp,notnull,default:true" json:"is_password_otp"`
	CreatedAt     time.Time   `bun:"created_at,notnull" json:"created_at,omitempty"`
	CustomUser    *CustomUser `bun:"rel:belongs-to,join:custom_user_id=id" json:"custom_user,omitempty"`
}

// BeforeInsert hook executed before database insert operation.
func (p *PedagogicalSpecialist) BeforeInsert(db *bun.DB) error {
	p.CreatedAt = time.Now()
	return p.Validate()
}

// Validate validates PedagogicalSpecialist struct and returns validation errors.
func (p *PedagogicalSpecialist) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Role, validation.Required, validation.Length(1, 50)),
		validation.Field(&p.CustomUserID, validation.Required),
		validation.Field(&p.UserID, validation.Required),
	)
}
