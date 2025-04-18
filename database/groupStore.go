package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dhax/go-base/models"
	"github.com/uptrace/bun"
)

// GroupStore implements database operations for group management
type GroupStore struct {
	db *bun.DB
}

// NewGroupStore returns a GroupStore
func NewGroupStore(db *bun.DB) *GroupStore {
	return &GroupStore{
		db: db,
	}
}

// GetGroupByID retrieves a Group by ID with related data
func (s *GroupStore) GetGroupByID(ctx context.Context, id int64) (*models.Group, error) {
	group := new(models.Group)
	err := s.db.NewSelect().
		Model(group).
		Relation("Room").
		Relation("Representative").
		Relation("Representative.CustomUser").
		Relation("Supervisors").
		Relation("Supervisors.CustomUser").
		Where("\"group\".id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return group, nil
}

// CreateGroup creates a new Group with supervisors
func (s *GroupStore) CreateGroup(ctx context.Context, group *models.Group, supervisorIDs []int64) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create the group
	_, err = tx.NewInsert().
		Model(group).
		Exec(ctx)

	if err != nil {
		return err
	}

	// Add supervisors if provided
	if len(supervisorIDs) > 0 {
		for _, specialistID := range supervisorIDs {
			groupSupervisor := &models.GroupSupervisor{
				GroupID:      group.ID,
				SpecialistID: specialistID,
			}

			_, err = tx.NewInsert().
				Model(groupSupervisor).
				Exec(ctx)

			if err != nil {
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// UpdateGroup updates an existing Group
func (s *GroupStore) UpdateGroup(ctx context.Context, group *models.Group) error {
	_, err := s.db.NewUpdate().
		Model(group).
		Column("name", "room_id", "representative_id", "modified_at").
		WherePK().
		Exec(ctx)

	return err
}

// DeleteGroup deletes a Group
func (s *GroupStore) DeleteGroup(ctx context.Context, id int64) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if the group has students
	hasStudents, err := tx.NewSelect().
		Model((*models.Student)(nil)).
		Where("group_id = ?", id).
		Exists(ctx)

	if err != nil {
		return err
	}

	if hasStudents {
		return errors.New("cannot delete group with students")
	}

	// Delete group supervisors
	_, err = tx.NewDelete().
		Model((*models.GroupSupervisor)(nil)).
		Where("group_id = ?", id).
		Exec(ctx)

	if err != nil {
		return err
	}

	// Delete combined group relationships
	_, err = tx.NewDelete().
		Model((*models.CombinedGroupGroup)(nil)).
		Where("group_id = ?", id).
		Exec(ctx)

	if err != nil {
		return err
	}

	// Delete group
	_, err = tx.NewDelete().
		Model((*models.Group)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// ListGroups returns a list of all Groups with optional filtering
func (s *GroupStore) ListGroups(ctx context.Context, filters map[string]interface{}) ([]models.Group, error) {
	var groups []models.Group

	query := s.db.NewSelect().
		Model(&groups).
		Relation("Room").
		Relation("Representative").
		Relation("Representative.CustomUser").
		Relation("Supervisors").
		Relation("Supervisors.CustomUser")

	// Apply filters
	if supervisorID, ok := filters["supervisor_id"].(int64); ok && supervisorID > 0 {
		query = query.Join("JOIN group_supervisors gs ON gs.group_id = \"group\".id").
			Where("gs.specialist_id = ?", supervisorID)
	}

	if searchTerm, ok := filters["search"].(string); ok && searchTerm != "" {
		query = query.Where("\"group\".name ILIKE ?", "%"+searchTerm+"%")
	}

	err := query.OrderExpr("\"group\".name ASC").
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return groups, nil
}

// UpdateGroupSupervisors updates the supervisors for a group
func (s *GroupStore) UpdateGroupSupervisors(ctx context.Context, groupID int64, supervisorIDs []int64) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete existing supervisors
	_, err = tx.NewDelete().
		Model((*models.GroupSupervisor)(nil)).
		Where("group_id = ?", groupID).
		Exec(ctx)

	if err != nil {
		return err
	}

	// Add new supervisors
	for _, specialistID := range supervisorIDs {
		groupSupervisor := &models.GroupSupervisor{
			GroupID:      groupID,
			SpecialistID: specialistID,
		}

		_, err = tx.NewInsert().
			Model(groupSupervisor).
			Exec(ctx)

		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// CreateCombinedGroup creates a new CombinedGroup
func (s *GroupStore) CreateCombinedGroup(ctx context.Context, combinedGroup *models.CombinedGroup, groupIDs []int64, specialistIDs []int64) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create the combined group
	_, err = tx.NewInsert().
		Model(combinedGroup).
		Exec(ctx)

	if err != nil {
		return err
	}

	// Add groups
	for _, groupID := range groupIDs {
		combinedGroupGroup := &models.CombinedGroupGroup{
			CombinedGroupID: combinedGroup.ID,
			GroupID:         groupID,
		}

		_, err = tx.NewInsert().
			Model(combinedGroupGroup).
			Exec(ctx)

		if err != nil {
			return err
		}
	}

	// Add specialists
	for _, specialistID := range specialistIDs {
		combinedGroupSpecialist := &models.CombinedGroupSpecialist{
			CombinedGroupID: combinedGroup.ID,
			SpecialistID:    specialistID,
		}

		_, err = tx.NewInsert().
			Model(combinedGroupSpecialist).
			Exec(ctx)

		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// GetCombinedGroupByID retrieves a CombinedGroup by ID with related data
func (s *GroupStore) GetCombinedGroupByID(ctx context.Context, id int64) (*models.CombinedGroup, error) {
	combinedGroup := new(models.CombinedGroup)
	err := s.db.NewSelect().
		Model(combinedGroup).
		Relation("SpecificGroup").
		Relation("Groups").
		Relation("AccessSpecialists").
		Relation("AccessSpecialists.CustomUser").
		Where("combined_group.id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return combinedGroup, nil
}

// ListCombinedGroups returns a list of all CombinedGroups
func (s *GroupStore) ListCombinedGroups(ctx context.Context) ([]models.CombinedGroup, error) {
	var combinedGroups []models.CombinedGroup

	err := s.db.NewSelect().
		Model(&combinedGroups).
		Relation("SpecificGroup").
		Relation("Groups").
		Relation("AccessSpecialists").
		Relation("AccessSpecialists.CustomUser").
		OrderExpr("combined_group.name ASC").
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return combinedGroups, nil
}

// MergeRooms merges two rooms and creates a combined group
func (s *GroupStore) MergeRooms(ctx context.Context, sourceRoomID, targetRoomID int64) (*models.CombinedGroup, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Get source room and group
	var sourceGroup models.Group
	err = tx.NewSelect().
		Model(&sourceGroup).
		Where("room_id = ?", sourceRoomID).
		Scan(ctx)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, errors.New("no group found for source room")
	}

	// Get target room and group
	var targetGroup models.Group
	err = tx.NewSelect().
		Model(&targetGroup).
		Where("room_id = ?", targetRoomID).
		Scan(ctx)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, errors.New("no group found for target room")
	}

	// Create a combined group with both groups
	combinedGroupName := fmt.Sprintf("%s + %s", sourceGroup.Name, targetGroup.Name)
	combinedGroup := &models.CombinedGroup{
		Name:         combinedGroupName,
		IsActive:     true,
		AccessPolicy: "all", // All supervisors from both groups have access
	}

	_, err = tx.NewInsert().
		Model(combinedGroup).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	// Add both groups to the combined group
	sourceGroupLink := &models.CombinedGroupGroup{
		CombinedGroupID: combinedGroup.ID,
		GroupID:         sourceGroup.ID,
	}

	_, err = tx.NewInsert().
		Model(sourceGroupLink).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	targetGroupLink := &models.CombinedGroupGroup{
		CombinedGroupID: combinedGroup.ID,
		GroupID:         targetGroup.ID,
	}

	_, err = tx.NewInsert().
		Model(targetGroupLink).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	// Get supervisors for both groups
	var sourceSupervisors []models.GroupSupervisor
	err = tx.NewSelect().
		Model(&sourceSupervisors).
		Where("group_id = ?", sourceGroup.ID).
		Scan(ctx)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	var targetSupervisors []models.GroupSupervisor
	err = tx.NewSelect().
		Model(&targetSupervisors).
		Where("group_id = ?", targetGroup.ID).
		Scan(ctx)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Add all supervisors to the combined group
	addedSpecialists := make(map[int64]bool)

	for _, supervisor := range sourceSupervisors {
		if !addedSpecialists[supervisor.SpecialistID] {
			combinedGroupSpecialist := &models.CombinedGroupSpecialist{
				CombinedGroupID: combinedGroup.ID,
				SpecialistID:    supervisor.SpecialistID,
			}

			_, err = tx.NewInsert().
				Model(combinedGroupSpecialist).
				Exec(ctx)

			if err != nil {
				return nil, err
			}

			addedSpecialists[supervisor.SpecialistID] = true
		}
	}

	for _, supervisor := range targetSupervisors {
		if !addedSpecialists[supervisor.SpecialistID] {
			combinedGroupSpecialist := &models.CombinedGroupSpecialist{
				CombinedGroupID: combinedGroup.ID,
				SpecialistID:    supervisor.SpecialistID,
			}

			_, err = tx.NewInsert().
				Model(combinedGroupSpecialist).
				Exec(ctx)

			if err != nil {
				return nil, err
			}

			addedSpecialists[supervisor.SpecialistID] = true
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Return the created combined group with all relations
	return s.GetCombinedGroupByID(ctx, combinedGroup.ID)
}
