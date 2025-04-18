package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRoom_Validate(t *testing.T) {
	tests := []struct {
		name    string
		room    *Room
		wantErr bool
	}{
		{
			name: "valid room",
			room: &Room{
				RoomName: "Classroom 101",
				Building: "Main Building",
				Floor:    1,
				Capacity: 30,
				Category: "Classroom",
				Color:    "#FFFFFF",
			},
			wantErr: false,
		},
		{
			name: "missing room name",
			room: &Room{
				RoomName: "",
				Building: "Main Building",
				Floor:    1,
				Capacity: 30,
				Category: "Classroom",
				Color:    "#FFFFFF",
			},
			wantErr: true,
		},
		{
			name: "invalid capacity",
			room: &Room{
				RoomName: "Classroom 101",
				Building: "Main Building",
				Floor:    1,
				Capacity: 0, // Invalid: must be at least 1
				Category: "Classroom",
				Color:    "#FFFFFF",
			},
			wantErr: true,
		},
		{
			name: "invalid color",
			room: &Room{
				RoomName: "Classroom 101",
				Building: "Main Building",
				Floor:    1,
				Capacity: 30,
				Category: "Classroom",
				Color:    "not-a-hex-color", // Invalid hex color
			},
			wantErr: true,
		},
		{
			name: "missing category",
			room: &Room{
				RoomName: "Classroom 101",
				Building: "Main Building",
				Floor:    1,
				Capacity: 30,
				Category: "", // Required
				Color:    "#FFFFFF",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.room.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRoom_BeforeInsert(t *testing.T) {
	room := &Room{
		RoomName: "Test Room",
		Building: "Test Building",
		Floor:    1,
		Capacity: 25,
		Category: "Test",
		Color:    "#FF0000",
	}

	// Timestamps should be zero before BeforeInsert
	assert.True(t, room.CreatedAt.IsZero())
	assert.True(t, room.UpdatedAt.IsZero())

	// Call BeforeInsert (we'll pass nil as DB since we're not actually inserting)
	room.BeforeInsert(nil)

	// Check that timestamps are set
	assert.False(t, room.CreatedAt.IsZero())
	assert.False(t, room.UpdatedAt.IsZero())

	// CreatedAt and UpdatedAt should be very close to current time
	now := time.Now()
	assert.WithinDuration(t, now, room.CreatedAt, 2*time.Second)
	assert.WithinDuration(t, now, room.UpdatedAt, 2*time.Second)
}

func TestRoom_BeforeUpdate(t *testing.T) {
	room := &Room{
		RoomName:  "Test Room",
		Building:  "Test Building",
		Floor:     1,
		Capacity:  25,
		Category:  "Test",
		Color:     "#FF0000",
		CreatedAt: time.Now().Add(-24 * time.Hour), // Set CreatedAt to yesterday
	}

	// UpdatedAt should be zero before BeforeUpdate
	assert.True(t, room.UpdatedAt.IsZero())

	// Call BeforeUpdate (we'll pass nil as DB since we're not actually updating)
	room.BeforeUpdate(nil)

	// Check that UpdatedAt is set but CreatedAt hasn't changed
	assert.False(t, room.UpdatedAt.IsZero())
	assert.WithinDuration(t, time.Now(), room.UpdatedAt, 2*time.Second)
	assert.WithinDuration(t, time.Now().Add(-24*time.Hour), room.CreatedAt, 2*time.Second)
}
