package database

import (
	"context"
	"time"

	"github.com/dhax/go-base/models"
	"github.com/uptrace/bun"
)

// CustomUserStore implements database operations for custom user management.
type CustomUserStore struct {
	db *bun.DB
}

// NewCustomUserStore returns a new CustomUserStore instance.
func NewCustomUserStore(db *bun.DB) *CustomUserStore {
	return &CustomUserStore{
		db: db,
	}
}

// Create creates a new custom user.
func (s *CustomUserStore) Create(ctx context.Context, user *models.CustomUser) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := s.db.NewInsert().
		Model(user).
		Exec(ctx)

	return err
}

// Get returns a custom user by ID.
func (s *CustomUserStore) Get(ctx context.Context, id int64) (*models.CustomUser, error) {
	user := &models.CustomUser{}
	err := s.db.NewSelect().
		Model(user).
		Where("id = ?", id).
		Scan(ctx)

	return user, err
}

// GetByAccountID returns a custom user by account ID.
func (s *CustomUserStore) GetByAccountID(ctx context.Context, accountID int64) (*models.CustomUser, error) {
	user := &models.CustomUser{}
	err := s.db.NewSelect().
		Model(user).
		Where("account_id = ?", accountID).
		Scan(ctx)

	return user, err
}

// GetByTagID returns a custom user by RFID tag ID.
func (s *CustomUserStore) GetByTagID(ctx context.Context, tagID string) (*models.CustomUser, error) {
	user := &models.CustomUser{}
	err := s.db.NewSelect().
		Model(user).
		Where("tag_id = ?", tagID).
		Scan(ctx)

	return user, err
}

// Update updates a custom user.
func (s *CustomUserStore) Update(ctx context.Context, user *models.CustomUser) error {
	user.UpdatedAt = time.Now()

	_, err := s.db.NewUpdate().
		Model(user).
		Where("id = ?", user.ID).
		Exec(ctx)

	return err
}

// Delete deletes a custom user.
func (s *CustomUserStore) Delete(ctx context.Context, id int64) error {
	_, err := s.db.NewDelete().
		Model((*models.CustomUser)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	return err
}

// List returns a list of custom users with optional filtering.
func (s *CustomUserStore) List(ctx context.Context, filters map[string]interface{}) ([]models.CustomUser, error) {
	var users []models.CustomUser

	query := s.db.NewSelect().
		Model(&users)

	// Apply filters if provided
	if search, ok := filters["search"].(string); ok && search != "" {
		query = query.Where("first_name ILIKE ?", "%"+search+"%").
			WhereOr("second_name ILIKE ?", "%"+search+"%").
			WhereOr("tag_id ILIKE ?", "%"+search+"%")
	}

	err := query.Order("first_name ASC, second_name ASC").
		Scan(ctx)

	return users, err
}

// UpdateTagID updates a user's RFID tag ID.
func (s *CustomUserStore) UpdateTagID(ctx context.Context, userID int64, tagID string) error {
	_, err := s.db.NewUpdate().
		Model((*models.CustomUser)(nil)).
		Set("tag_id = ?", tagID).
		Set("updated_at = ?", time.Now()).
		Where("id = ?", userID).
		Exec(ctx)

	return err
}
