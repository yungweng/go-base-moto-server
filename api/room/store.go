package room

import (
	"context"
	"fmt"
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
	// 1. Query all active RoomOccupancy entries
	var occupancies []RoomOccupancy
	err := s.db.NewSelect().
		Model(&occupancies).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	// 2. For each occupancy, construct a RoomOccupancyDetail
	var details []RoomOccupancyDetail
	for _, occupancy := range occupancies {
		// Get room
		room := new(models.Room)
		err := s.db.NewSelect().
			Model(room).
			Where("id = ?", occupancy.RoomID).
			Scan(ctx)
		if err != nil {
			continue // Skip if we can't find the room
		}

		// Get timespan
		timespan := new(models.Timespan)
		err = s.db.NewSelect().
			Model(timespan).
			Where("id = ?", occupancy.TimespanID).
			Scan(ctx)
		if err != nil {
			continue // Skip if we can't find the timespan
		}

		// Only include active timespans (no end time or end time in the future)
		if timespan.EndTime != nil && timespan.EndTime.Before(time.Now()) {
			continue
		}

		// Get supervisors
		var supervisors []SupervisorInfo
		err = s.db.NewSelect().
			Table("room_occupancy_supervisors").
			Column("specialist_id AS id").
			Where("room_occupancy_id = ?", occupancy.ID).
			Scan(ctx, &supervisors)
		if err != nil {
			supervisors = []SupervisorInfo{} // Continue without supervisors
		}

		// Construct RoomOccupancyDetail
		detail := RoomOccupancyDetail{
			Room: RoomInfo{
				RoomName: room.RoomName,
				Floor:    room.Floor,
				Capacity: room.Capacity,
			},
			Supervisor: supervisors,
			Timespan: TimespanInfo{
				StartTime: timespan.StartTime.Format("15:04"),
			},
		}

		if timespan.EndTime != nil {
			endTimeStr := timespan.EndTime.Format("15:04")
			detail.Timespan.EndTime = endTimeStr
		}

		details = append(details, detail)
	}

	return details, nil
}

// GetRoomOccupancyByID returns room occupancy details by ID
func (s *roomStore) GetRoomOccupancyByID(ctx context.Context, id int64) (*RoomOccupancyDetail, error) {
	// 1. Query RoomOccupancy by ID
	occupancy := new(RoomOccupancy)
	err := s.db.NewSelect().
		Model(occupancy).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	// 2. Get Room
	room := new(models.Room)
	err = s.db.NewSelect().
		Model(room).
		Where("id = ?", occupancy.RoomID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	// 3. Get Timespan
	timespan := new(models.Timespan)
	err = s.db.NewSelect().
		Model(timespan).
		Where("id = ?", occupancy.TimespanID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	// 4. Get Supervisors
	var supervisors []SupervisorInfo
	err = s.db.NewSelect().
		Table("room_occupancy_supervisors").
		Column("specialist_id AS id").
		Where("room_occupancy_id = ?", occupancy.ID).
		Scan(ctx, &supervisors)
	if err != nil {
		supervisors = []SupervisorInfo{} // Continue without supervisors
	}

	// 5. Construct RoomOccupancyDetail
	detail := &RoomOccupancyDetail{
		Room: RoomInfo{
			RoomName: room.RoomName,
			Floor:    room.Floor,
			Capacity: room.Capacity,
		},
		Supervisor: supervisors,
		Timespan: TimespanInfo{
			StartTime: timespan.StartTime.Format("15:04"),
		},
	}

	if timespan.EndTime != nil {
		endTimeStr := timespan.EndTime.Format("15:04")
		detail.Timespan.EndTime = endTimeStr
	}

	// AG info would be added when AG model is implemented

	return detail, nil
}

// GetCurrentRoomOccupancy returns current occupancy for a room
func (s *roomStore) GetCurrentRoomOccupancy(ctx context.Context, roomID int64) (*RoomOccupancyDetail, error) {
	// 1. Query Room to make sure it exists
	room := new(models.Room)
	err := s.db.NewSelect().
		Model(room).
		Where("id = ?", roomID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	// 2. Query RoomOccupancy by roomID
	occupancy := new(RoomOccupancy)
	err = s.db.NewSelect().
		Model(occupancy).
		Where("room_id = ?", roomID).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("room is not currently occupied")
	}

	// 3. Get Timespan for the occupancy
	timespan := new(models.Timespan)
	err = s.db.NewSelect().
		Model(timespan).
		Where("id = ?", occupancy.TimespanID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	// Only include active timespans (no end time or end time in the future)
	if timespan.EndTime != nil && timespan.EndTime.Before(time.Now()) {
		return nil, fmt.Errorf("room is not currently occupied")
	}

	// 4. Get supervisors for the occupancy
	var supervisors []SupervisorInfo
	err = s.db.NewSelect().
		Table("room_occupancy_supervisors").
		Column("specialist_id AS id").
		Where("room_occupancy_id = ?", occupancy.ID).
		Scan(ctx, &supervisors)
	if err != nil {
		// Continue without supervisors if we can't fetch them
		supervisors = []SupervisorInfo{}
	}

	// For now, we'll return limited supervisor info as we don't have access to the specialist model yet

	// 5. Construct RoomOccupancyDetail object
	detail := &RoomOccupancyDetail{
		Room: RoomInfo{
			RoomName: room.RoomName,
			Floor:    room.Floor,
			Capacity: room.Capacity,
		},
		Supervisor: supervisors,
		Timespan: TimespanInfo{
			StartTime: timespan.StartTime.Format("15:04"),
		},
	}

	if timespan.EndTime != nil {
		endTimeStr := timespan.EndTime.Format("15:04")
		detail.Timespan.EndTime = endTimeStr
	}

	// We'll add AG info when AG model is implemented

	return detail, nil
}

// RegisterTablet registers a tablet to a room
func (s *roomStore) RegisterTablet(ctx context.Context, roomID int64, req *RegisterTabletRequest) (*RoomOccupancy, error) {
	// Start a transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 1. Check if room exists
	room := new(models.Room)
	err = tx.NewSelect().
		Model(room).
		Where("id = ?", roomID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	// 2. Check if tablet is already registered
	var existingOccupancy RoomOccupancy
	err = tx.NewSelect().
		Model(&existingOccupancy).
		Where("device_id = ?", req.DeviceID).
		Scan(ctx)
	if err == nil {
		// Tablet is already registered
		return nil, fmt.Errorf("tablet is already registered")
	}

	// 3. Create a timespan
	timespan := &models.Timespan{
		StartTime: time.Now(),
		CreatedAt: time.Now(),
	}
	_, err = tx.NewInsert().
		Model(timespan).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	// 4. Create a new AG if requested (placeholder for now)
	var agID int64
	if req.NewAg != nil {
		// AG creation logic would go here
		// For now, we'll just use the provided AgID if any
		if req.AgID != nil {
			agID = *req.AgID
		}
	} else if req.AgID != nil {
		agID = *req.AgID
	}

	// 5. Create RoomOccupancy entry
	occupancy := &RoomOccupancy{
		DeviceID:   req.DeviceID,
		RoomID:     roomID,
		TimespanID: timespan.ID,
		CreatedAt:  time.Now(),
	}

	if req.GroupID != nil {
		occupancy.GroupID = *req.GroupID
	}

	if agID != 0 {
		occupancy.AgID = agID
	}

	_, err = tx.NewInsert().
		Model(occupancy).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	// 6. Add supervisors to RoomOccupancySupervisor table
	for _, supervisorID := range req.Supervisors {
		supervisor := &RoomOccupancySupervisor{
			RoomOccupancyID: occupancy.ID,
			SpecialistID:    supervisorID,
			CreatedAt:       time.Now(),
		}
		_, err = tx.NewInsert().
			Model(supervisor).
			Exec(ctx)
		if err != nil {
			return nil, err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return occupancy, nil
}

// UnregisterTablet unregisters a tablet from a room
func (s *roomStore) UnregisterTablet(ctx context.Context, roomID int64, deviceID string) error {
	// Start a transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Find RoomOccupancy by roomID and deviceID
	occupancy := new(RoomOccupancy)
	err = tx.NewSelect().
		Model(occupancy).
		Where("room_id = ? AND device_id = ?", roomID, deviceID).
		Scan(ctx)
	if err != nil {
		return fmt.Errorf("tablet not registered to this room")
	}

	// Get the timespan to update its end time
	timespan := new(models.Timespan)
	err = tx.NewSelect().
		Model(timespan).
		Where("id = ?", occupancy.TimespanID).
		Scan(ctx)
	if err != nil {
		return err
	}

	// Update the timespan to mark the end time
	endTime := time.Now()
	timespan.EndTime = &endTime
	_, err = tx.NewUpdate().
		Model(timespan).
		Column("endtime").
		Where("id = ?", timespan.ID).
		Exec(ctx)
	if err != nil {
		return err
	}

	// 2. Delete related RoomOccupancySupervisor entries
	_, err = tx.NewDelete().
		Model((*RoomOccupancySupervisor)(nil)).
		Where("room_occupancy_id = ?", occupancy.ID).
		Exec(ctx)
	if err != nil {
		return err
	}

	// 3. Delete RoomOccupancy entry
	_, err = tx.NewDelete().
		Model(occupancy).
		WherePK().
		Exec(ctx)
	if err != nil {
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

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
