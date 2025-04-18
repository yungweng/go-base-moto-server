package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/uptrace/bun"
)

// Group represents a collection of students
type Group struct {
	bun.BaseModel `bun:"table:groups,alias:g"`

	ID               int64                  `bun:"id,pk,autoincrement" json:"id"`
	Name             string                 `bun:"name,notnull,unique" json:"name"`
	RoomID           *int64                 `bun:"room_id" json:"room_id,omitempty"`
	RepresentativeID *int64                 `bun:"representative_id" json:"representative_id,omitempty"`
	CreatedAt        time.Time              `bun:"created_at,notnull" json:"created_at,omitempty"`
	UpdatedAt        time.Time              `bun:"updated_at,notnull" json:"updated_at,omitempty"`
	Room             *Room                  `bun:"rel:belongs-to,join:room_id=id" json:"room,omitempty"`
	Representative   *PedagogicalSpecialist `bun:"rel:belongs-to,join:representative_id=id" json:"representative,omitempty"`
	Students         []Student              `bun:"rel:has-many,join:id=group_id" json:"students,omitempty"`
}

// BeforeInsert hook executed before database insert operation.
func (g *Group) BeforeInsert(db *bun.DB) error {
	now := time.Now()
	g.CreatedAt = now
	g.UpdatedAt = now
	return g.Validate()
}

// BeforeUpdate hook executed before database update operation.
func (g *Group) BeforeUpdate(db *bun.DB) error {
	g.UpdatedAt = time.Now()
	return g.Validate()
}

// Validate validates Group struct and returns validation errors.
func (g *Group) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.Name, validation.Required, validation.Length(1, 100)),
	)
}

// CombinedGroup represents a temporary grouping of multiple groups
type CombinedGroup struct {
	bun.BaseModel `bun:"table:combined_groups,alias:cg"`

	ID              int64     `bun:"id,pk,autoincrement" json:"id"`
	Name            string    `bun:"name,notnull,unique" json:"name"`
	IsActive        bool      `bun:"is_active,notnull,default:true" json:"is_active"`
	CreatedAt       time.Time `bun:"created_at,notnull" json:"created_at,omitempty"`
	ValidUntil      time.Time `bun:"valid_until" json:"valid_until,omitempty"`
	AccessPolicy    string    `bun:"access_policy,notnull" json:"access_policy"`
	SpecificGroupID *int64    `bun:"specific_group_id" json:"specific_group_id,omitempty"`
	SpecificGroup   *Group    `bun:"rel:belongs-to,join:specific_group_id=id" json:"specific_group,omitempty"`
	// We'll implement this relationship programmatically instead of using m2m
	Groups []Group `json:"groups,omitempty" bun:"-"`
}

// BeforeInsert hook executed before database insert operation.
func (cg *CombinedGroup) BeforeInsert(db *bun.DB) error {
	cg.CreatedAt = time.Now()
	return cg.Validate()
}

// Validate validates CombinedGroup struct and returns validation errors.
func (cg *CombinedGroup) Validate() error {
	return validation.ValidateStruct(cg,
		validation.Field(&cg.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&cg.AccessPolicy, validation.Required, validation.In("all", "first", "specific", "manual")),
	)
}
