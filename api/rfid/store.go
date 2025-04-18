package rfid

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

// RFIDStore defines database operations for RFID tag management
type RFIDStore interface {
	SaveTag(ctx context.Context, tagID, readerID string) (*Tag, error)
	GetAllTags(ctx context.Context) ([]Tag, error)
	GetTagStats(ctx context.Context) (int, error)
	SaveTauriTags(ctx context.Context, deviceID string, tags []SyncTag) error
}

type rfidStore struct {
	db *bun.DB
}

// NewRFIDStore returns a new RFIDStore implementation
func NewRFIDStore(db *bun.DB) RFIDStore {
	return &rfidStore{db: db}
}

// SaveTag saves an RFID tag read to the database
func (s *rfidStore) SaveTag(ctx context.Context, tagID, readerID string) (*Tag, error) {
	now := time.Now()
	tag := &Tag{
		TagID:     tagID,
		ReaderID:  readerID,
		ReadAt:    now,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err := s.db.NewInsert().
		Model(tag).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	return tag, nil
}

// GetAllTags returns all stored RFID tags
func (s *rfidStore) GetAllTags(ctx context.Context) ([]Tag, error) {
	var tags []Tag
	err := s.db.NewSelect().
		Model(&tags).
		Order("read_at DESC").
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return tags, nil
}

// GetTagStats returns statistics about RFID tags
func (s *rfidStore) GetTagStats(ctx context.Context) (int, error) {
	count, err := s.db.NewSelect().
		Model((*Tag)(nil)).
		Count(ctx)

	if err != nil {
		return 0, err
	}

	return count, nil
}

// SaveTauriTags saves a batch of tags from the Tauri app
func (s *rfidStore) SaveTauriTags(ctx context.Context, deviceID string, syncTags []SyncTag) error {
	if len(syncTags) == 0 {
		return nil
	}

	tags := make([]Tag, len(syncTags))
	now := time.Now()

	for i, syncTag := range syncTags {
		tags[i] = Tag{
			TagID:     syncTag.TagID,
			ReaderID:  syncTag.ReaderID,
			ReadAt:    syncTag.LocalReadAt,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	_, err := s.db.NewInsert().
		Model(&tags).
		Exec(ctx)

	return err
}