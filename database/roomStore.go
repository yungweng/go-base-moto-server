package database

import (
	"context"
	"time"

	"github.com/dhax/go-base/models"
	"github.com/uptrace/bun"
)

// RoomStore implements database operations for room management.
type RoomStore struct {
	db *bun.DB
}

// NewRoomStore returns a new RoomStore instance.
func NewRoomStore(db *bun.DB) *RoomStore {
	return &RoomStore{
		db: db,
	}
}

// Create creates a new room.
func (s *RoomStore) Create(ctx context.Context, room *models.Room) error {
	now := time.Now()
	room.CreatedAt = now
	room.UpdatedAt = now

	_, err := s.db.NewInsert().
		Model(room).
		Exec(ctx)

	return err
}

// Get returns a room by ID.
func (s *RoomStore) Get(ctx context.Context, id int64) (*models.Room, error) {
	room := &models.Room{}
	err := s.db.NewSelect().
		Model(room).
		Where("id = ?", id).
		Scan(ctx)

	return room, err
}

// Update updates a room.
func (s *RoomStore) Update(ctx context.Context, room *models.Room) error {
	room.UpdatedAt = time.Now()

	_, err := s.db.NewUpdate().
		Model(room).
		Where("id = ?", room.ID).
		Exec(ctx)

	return err
}

// Delete deletes a room.
func (s *RoomStore) Delete(ctx context.Context, id int64) error {
	_, err := s.db.NewDelete().
		Model((*models.Room)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	return err
}

// List returns a list of rooms with optional filtering.
func (s *RoomStore) List(ctx context.Context, filters map[string]interface{}) ([]models.Room, error) {
	var rooms []models.Room

	query := s.db.NewSelect().
		Model(&rooms)

	// Apply filters if provided
	if category, ok := filters["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}

	if building, ok := filters["building"].(string); ok && building != "" {
		query = query.Where("building = ?", building)
	}

	if floor, ok := filters["floor"].(int); ok {
		query = query.Where("floor = ?", floor)
	}

	if search, ok := filters["search"].(string); ok && search != "" {
		query = query.Where("room_name ILIKE ?", "%"+search+"%").
			WhereOr("category ILIKE ?", "%"+search+"%").
			WhereOr("building ILIKE ?", "%"+search+"%")
	}

	err := query.Order("room_name ASC").
		Scan(ctx)

	return rooms, err
}

// GroupByCategory returns rooms grouped by their categories.
func (s *RoomStore) GroupByCategory(ctx context.Context) (map[string][]models.Room, error) {
	var rooms []models.Room

	err := s.db.NewSelect().
		Model(&rooms).
		Order("category ASC, room_name ASC").
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	// Group rooms by category
	grouped := make(map[string][]models.Room)
	for _, room := range rooms {
		grouped[room.Category] = append(grouped[room.Category], room)
	}

	return grouped, nil
}
