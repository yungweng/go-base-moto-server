package rfid

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dhax/go-base/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestHandleRoomEntry tests the room entry endpoint
func TestHandleRoomEntry(t *testing.T) {
	// Setup
	mockRFIDStore := new(MockRFIDStore)
	mockUserStore := new(MockUserStore)
	mockStudentStore := new(MockStudentStore)

	api := &API{
		store:        mockRFIDStore,
		userStore:    mockUserStore,
		studentStore: mockStudentStore,
	}

	// Mock data
	tagID := "TAG123456"
	readerID := "READER001"
	roomID := int64(42)
	now := time.Now()

	mockTag := &Tag{
		ID:        1,
		TagID:     tagID,
		ReaderID:  readerID,
		ReadAt:    now,
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockUser := &models.CustomUser{
		ID:         24,
		FirstName:  "John",
		SecondName: "Doe",
		TagID:      &tagID,
	}

	mockStudent := &models.Student{
		ID:           123,
		SchoolClass:  "1A",
		CustomUserID: 24,
	}

	// Set expectations
	mockRFIDStore.On("SaveTag", mock.Anything, tagID, readerID).Return(mockTag, nil)
	mockUserStore.On("GetCustomUserByTagID", mock.Anything, tagID).Return(mockUser, nil)
	mockStudentStore.On("GetStudentByCustomUserID", mock.Anything, int64(24)).Return(mockStudent, nil)
	mockRFIDStore.On("RecordRoomEntry", mock.Anything, int64(123), roomID).Return(nil)

	// Expect the location update with correct parameters
	mockStudentStore.On("UpdateStudentLocation", mock.Anything, int64(123), mock.MatchedBy(func(locations map[string]bool) bool {
		return locations["in_house"] == true && locations["wc"] == false && locations["school_yard"] == false
	})).Return(nil)

	// Mock for GetRoomOccupancy
	mockOccupancy := &RoomOccupancyData{
		RoomID:       roomID,
		RoomName:     "Test Room",
		Capacity:     30,
		StudentCount: 1,
		Students: []RoomOccupancyStudent{
			{
				ID:        123,
				Name:      "John Doe",
				EnteredAt: now,
			},
		},
	}
	mockRFIDStore.On("GetRoomOccupancy", mock.Anything, roomID).Return(mockOccupancy, nil)

	// Create request
	payload := `{"tag_id":"TAG123456","room_id":42,"reader_id":"READER001"}`
	req := httptest.NewRequest("POST", "/room-entry", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	api.handleRoomEntry(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response OccupancyResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Student entered room successfully", response.Message)
	assert.Equal(t, int64(123), response.StudentID)
	assert.Equal(t, roomID, response.RoomID)
	assert.Equal(t, 1, response.StudentCount)

	// Verify all expectations were met
	mockRFIDStore.AssertExpectations(t)
	mockUserStore.AssertExpectations(t)
	mockStudentStore.AssertExpectations(t)
}

// Mock RoomOccupancy methods
func (m *MockRFIDStore) RecordRoomEntry(ctx context.Context, studentID, roomID int64) error {
	args := m.Called(ctx, studentID, roomID)
	return args.Error(0)
}

func (m *MockRFIDStore) RecordRoomExit(ctx context.Context, studentID, roomID int64) error {
	args := m.Called(ctx, studentID, roomID)
	return args.Error(0)
}

func (m *MockRFIDStore) GetRoomOccupancy(ctx context.Context, roomID int64) (*RoomOccupancyData, error) {
	args := m.Called(ctx, roomID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RoomOccupancyData), args.Error(1)
}

func (m *MockRFIDStore) GetCurrentRooms(ctx context.Context) ([]RoomOccupancyData, error) {
	args := m.Called(ctx)
	return args.Get(0).([]RoomOccupancyData), args.Error(1)
}

// Mock StudentStore's ListStudents method
func (m *MockStudentStore) ListStudents(ctx context.Context, filters map[string]interface{}) ([]models.Student, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]models.Student), args.Error(1)
}
