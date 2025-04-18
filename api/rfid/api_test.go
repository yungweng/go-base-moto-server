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

// Mock RFIDStore
type MockRFIDStore struct {
	mock.Mock
}

func (m *MockRFIDStore) SaveTag(ctx context.Context, tagID, readerID string) (*Tag, error) {
	args := m.Called(ctx, tagID, readerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Tag), args.Error(1)
}

func (m *MockRFIDStore) GetAllTags(ctx context.Context) ([]Tag, error) {
	args := m.Called(ctx)
	return args.Get(0).([]Tag), args.Error(1)
}

func (m *MockRFIDStore) GetTagStats(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockRFIDStore) SaveTauriTags(ctx context.Context, deviceID string, tags []SyncTag) error {
	args := m.Called(ctx, deviceID, tags)
	return args.Error(0)
}

// Mock UserStore
type MockUserStore struct {
	mock.Mock
}

func (m *MockUserStore) GetCustomUserByTagID(ctx context.Context, tagID string) (*models.CustomUser, error) {
	args := m.Called(ctx, tagID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CustomUser), args.Error(1)
}

// Mock StudentStore
type MockStudentStore struct {
	mock.Mock
}

func (m *MockStudentStore) GetStudentByCustomUserID(ctx context.Context, customUserID int64) (*models.Student, error) {
	args := m.Called(ctx, customUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Student), args.Error(1)
}

func (m *MockStudentStore) UpdateStudentLocation(ctx context.Context, id int64, locations map[string]bool) error {
	args := m.Called(ctx, id, locations)
	return args.Error(0)
}

func TestHandleTagRead(t *testing.T) {
	// Setup
	mockStore := new(MockRFIDStore)
	api := &API{store: mockStore}

	// Mock data
	tagID := "ABCDEF123456"
	readerID := "READER001"
	now := time.Now()

	mockTag := &Tag{
		ID:        1,
		TagID:     tagID,
		ReaderID:  readerID,
		ReadAt:    now,
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockStore.On("SaveTag", mock.Anything, tagID, readerID).Return(mockTag, nil)

	// Create request
	payload := `{"tag_id":"ABCDEF123456","reader_id":"READER001"}`
	req := httptest.NewRequest("POST", "/tag", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	api.handleTagRead(w, req)

	// Assert response
	assert.Equal(t, http.StatusCreated, w.Code)

	var response Tag
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), response.ID)
	assert.Equal(t, tagID, response.TagID)
	assert.Equal(t, readerID, response.ReaderID)

	mockStore.AssertExpectations(t)
}

func TestHandleStudentTracking(t *testing.T) {
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
	tagID := "ABCDEF123456"
	readerID := "READER001"
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
		ID:         42,
		FirstName:  "John",
		SecondName: "Doe",
		TagID:      &tagID,
	}

	mockStudent := &models.Student{
		ID:           24,
		SchoolClass:  "1A",
		CustomUserID: 42,
	}

	// Set expectations
	mockRFIDStore.On("SaveTag", mock.Anything, tagID, readerID).Return(mockTag, nil)
	mockUserStore.On("GetCustomUserByTagID", mock.Anything, tagID).Return(mockUser, nil)
	mockStudentStore.On("GetStudentByCustomUserID", mock.Anything, int64(42)).Return(mockStudent, nil)

	// Expect the location update with correct parameters
	mockStudentStore.On("UpdateStudentLocation", mock.Anything, int64(24), mock.MatchedBy(func(locations map[string]bool) bool {
		return locations["in_house"] == true && locations["wc"] == false && locations["school_yard"] == false
	})).Return(nil)

	// Create request
	payload := `{"tag_id":"ABCDEF123456","reader_id":"READER001","location_type":"entry"}`
	req := httptest.NewRequest("POST", "/track-student", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	api.handleStudentTracking(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response StudentTrackingResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "John Doe", response.Name)
	assert.Equal(t, "in-house", response.Location)

	// Verify all expectations were met
	mockRFIDStore.AssertExpectations(t)
	mockUserStore.AssertExpectations(t)
	mockStudentStore.AssertExpectations(t)
}
