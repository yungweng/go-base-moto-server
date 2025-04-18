package models

import (
	"time"

	"github.com/uptrace/bun"
)

// CombinedGroupGroups is a junction table for M2M relationship between CombinedGroup and Group
type CombinedGroupGroups struct {
	bun.BaseModel `bun:"table:combined_group_groups,alias:cgg"`

	ID              int64          `bun:"id,pk,autoincrement" json:"id"`
	CombinedGroupID int64          `bun:"combined_group_id,notnull" json:"combined_group_id"`
	GroupID         int64          `bun:"group_id,notnull" json:"group_id"`
	CreatedAt       time.Time      `bun:"created_at,notnull" json:"created_at"`
	CombinedGroup   *CombinedGroup `bun:"rel:belongs-to,join:combined_group_id=id" json:"combined_group,omitempty"`
	Group           *Group         `bun:"rel:belongs-to,join:group_id=id" json:"group,omitempty"`
}

// BeforeInsert hook executed before database insert operation.
func (cgg *CombinedGroupGroups) BeforeInsert(db *bun.DB) error {
	cgg.CreatedAt = time.Now()
	return nil
}
