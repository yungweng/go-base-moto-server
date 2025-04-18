package room

import (
	"context"
	"time"

	"github.com/dhax/go-base/models"
	"github.com/uptrace/bun"
)

// RoomStore defines database operations for room management
type RoomStore interface {
	// Room operations
	GetRooms(ctx context.Context) ([]models.Room, error)
	GetRoomsByCategory(ctx context.Context, category string) ([]models.Room, error)
	GetRoomsByBuilding(ctx context.Context, building string) ([]models.Room, error)
	GetRoomsByFloor(ctx context.Context, floor int) ([]models.Room, error)
	GetRoomsByOccupied(ctx context.Context, occupied bool) ([]models.Room, error)
	GetRoomByID(ctx context.Context, id int64) (*models.Room, error)
	CreateRoom(ctx context.Context, room *models.Room) error
	UpdateRoom(ctx context.Context, room *models.Room) error
	DeleteRoom(ctx context.Context, id int64) error
	GetRoomsGroupedByCategory(ctx context.Context) (map[string][]models.Room, error)

	// Room occupancy operations
	GetAllRoomOccupancies(ctx context.Context) ([]RoomOccupancyDetail, error)
	GetRoomOccupancyByID(ctx context.Context, id int64) (*RoomOccupancyDetail, error)
	GetCurrentRoomOccupancy(ctx context.Context, roomID int64) (*RoomOccupancyDetail, error)
	RegisterTablet(ctx context.Context, roomID int64, req *RegisterTabletRequest) (*RoomOccupancy, error)
	UnregisterTablet(ctx context.Context, roomID int64, deviceID string) error
	AddSupervisorToRoomOccupancy(ctx context.Context, roomOccupancyID, supervisorID int64) error
}

type roomStore struct {
	db *bun.DB
}

// NewRoomStore returns a new RoomStore implementation
func NewRoomStore(db *bun.DB) RoomStore {
	return &roomStore{db: db}
}

// GetRooms returns all rooms
func (s *roomStore) GetRooms(ctx context.Context) ([]models.Room, error) {
	var rooms []models.Room
	err := s.db.NewSelect().
		Model(&rooms).
		Order("room_name ASC").
		Scan(ctx)

	return rooms, err
}

// GetRoomsByCategory returns rooms filtered by category
func (s *roomStore) GetRoomsByCategory(ctx context.Context, category string) ([]models.Room, error) {
	var rooms []models.Room
	err := s.db.NewSelect().
		Model(&rooms).
		Where("category = ?", category).
		Order("room_name ASC").
		Scan(ctx)

	return rooms, err
}

// GetRoomsByBuilding returns rooms filtered by building
func (s *roomStore) GetRoomsByBuilding(ctx context.Context, building string) ([]models.Room, error) {
	var rooms []models.Room
	err := s.db.NewSelect().
		Model(&rooms).
		Where("building = ?", building).
		Order("room_name ASC").
		Scan(ctx)

	return rooms, err
}

// GetRoomsByFloor returns rooms filtered by floor
func (s *roomStore) GetRoomsByFloor(ctx context.Context, floor int) ([]models.Room, error) {
	var rooms []models.Room
	err := s.db.NewSelect().
		Model(&rooms).
		Where("floor = ?", floor).
		Order("room_name ASC").
		Scan(ctx)

	return rooms, err
}

// GetRoomsByOccupied returns rooms filtered by occupancy status
func (s *roomStore) GetRoomsByOccupied(ctx context.Context, occupied bool) ([]models.Room, error) {
	var rooms []models.Room
	query := s.db.NewSelect().Model(&rooms)

	if occupied {
		// Join with RoomOccupancy to find occupied rooms
		query = query.Join("JOIN room_occupancies ro ON rooms.id = ro.room_id")
	} else {
		// Find rooms that don't have any occupancy entries
		query = query.Where("id NOT IN (SELECT room_id FROM room_occupancies)")
	}

	err := query.Order("room_name ASC").Scan(ctx)
	return rooms, err
}

// GetRoomByID returns a room by ID
func (s *roomStore) GetRoomByID(ctx context.Context, id int64) (*models.Room, error) {
	room := new(models.Room)
	err := s.db.NewSelect().
		Model(room).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return room, nil
}

// CreateRoom creates a new room
func (s *roomStore) CreateRoom(ctx context.Context, room *models.Room) error {
	room.CreatedAt = time.Now()
	room.ModifiedAt = time.Now()

	_, err := s.db.NewInsert().
		Model(room).
		Exec(ctx)

	return err
}

// UpdateRoom updates an existing room
func (s *roomStore) UpdateRoom(ctx context.Context, room *models.Room) error {
	room.ModifiedAt = time.Now()

	_, err := s.db.NewUpdate().
		Model(room).
		WherePK().
		Exec(ctx)

	return err
}

// DeleteRoom deletes a room by ID
func (s *roomStore) DeleteRoom(ctx context.Context, id int64) error {
	_, err := s.db.NewDelete().
		Model((*models.Room)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	return err
}

// GetRoomsGroupedByCategory returns rooms grouped by category
func (s *roomStore) GetRoomsGroupedByCategory(ctx context.Context) (map[string][]models.Room, error) {
	var rooms []models.Room
	err := s.db.NewSelect().
		Model(&rooms).
		Order("category ASC, room_name ASC").
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	// Group rooms by category
	groupedRooms := make(map[string][]models.Room)
	for _, room := range rooms {
		groupedRooms[room.Category] = append(groupedRooms[room.Category], room)
	}

	return groupedRooms, nil
}

// GetAllRoomOccupancies returns all room occupancies with details
func (s *roomStore) GetAllRoomOccupancies(ctx context.Context) ([]RoomOccupancyDetail, error) {
	// This is a placeholder - in a real implementation, you'd:
	// 1. Query RoomOccupancy entries
	// 2. For each entry, get Room, AG (if applicable), and Supervisors
	// 3. Construct RoomOccupancyDetail objects
	// This requires that AG, Timespan, and Specialist models are already implemented

	var details []RoomOccupancyDetail
	// Implementation would go here once other models are available
	return details, nil
}

// GetRoomOccupancyByID returns room occupancy details by ID
func (s *roomStore) GetRoomOccupancyByID(ctx context.Context, id int64) (*RoomOccupancyDetail, error) {
	// This is a placeholder - real implementation would:
	// 1. Query RoomOccupancy by ID
	// 2. Get related Room, AG (if applicable), and Supervisors
	// 3. Construct a RoomOccupancyDetail object

	// Implementation would go here once other models are available
	return nil, nil
}

// GetCurrentRoomOccupancy returns current occupancy for a room
func (s *roomStore) GetCurrentRoomOccupancy(ctx context.Context, roomID int64) (*RoomOccupancyDetail, error) {
	// This is a placeholder - real implementation would:
	// 1. Query RoomOccupancy by roomID
	// 2. Get related Room, AG (if applicable), and Supervisors
	// 3. Construct a RoomOccupancyDetail object

	// Implementation would go here once other models are available
	return nil, nil
}

// RegisterTablet registers a tablet to a room
func (s *roomStore) RegisterTablet(ctx context.Context, roomID int64, req *RegisterTabletRequest) (*RoomOccupancy, error) {
	// This is a placeholder - real implementation would:
	// 1. Check if room exists
	// 2. Check if tablet is already registered
	// 3. Create a timespan if needed
	// 4. Create a new AG if requested
	// 5. Create RoomOccupancy entry
	// 6. Add supervisors to RoomOccupancySupervisor table

	// Implementation would go here once other models are available
	return nil, nil
}

// UnregisterTablet unregisters a tablet from a room
func (s *roomStore) UnregisterTablet(ctx context.Context, roomID int64, deviceID string) error {
	// This is a placeholder - real implementation would:
	// 1. Find RoomOccupancy by roomID and deviceID
	// 2. Delete related RoomOccupancySupervisor entries
	// 3. Delete RoomOccupancy entry

	// Implementation would go here once other models are available
	return nil
}

// AddSupervisorToRoomOccupancy adds a supervisor to a room occupancy
func (s *roomStore) AddSupervisorToRoomOccupancy(ctx context.Context, roomOccupancyID, supervisorID int64) error {
	supervisor := &RoomOccupancySupervisor{
		RoomOccupancyID: roomOccupancyID,
		SpecialistID:    supervisorID,
		CreatedAt:       time.Now(),
	}

	_, err := s.db.NewInsert().
		Model(supervisor).
		Exec(ctx)

	return err
}
