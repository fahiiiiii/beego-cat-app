package controllers

import (
	"errors" // Import errors package
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"beego-cat-app/models" // Import the models package
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock dependencies
type MockBackendService struct {
	mock.Mock
}

func (m *MockBackendService) FetchCatImages() ([]models.CatImage, error) {
	args := m.Called()
	return args.Get(0).([]models.CatImage), args.Error(1)
}

// Define NewCatController for tests
func NewCatController(mockService *MockBackendService) *CatController {
	return &CatController{} // Modify as necessary to return a valid controller instance
}

// Test function to test the "GetImages" method in the controller
func TestGetImages(t *testing.T) {
	mockService := new(MockBackendService)

	// Mocking successful response
	mockService.On("FetchCatImages").Return([]models.CatImage{
		{ID: "1", URL: "http://example.com/cat1.jpg"},
		{ID: "2", URL: "http://example.com/cat2.jpg"},
	}, nil)

	// Creating a mock controller to test the method
	ctrl := NewCatController(mockService)
	req := httptest.NewRequest(http.MethodGet, "/api/getImages", nil)
	rec := httptest.NewRecorder()

	// Simulate the controller's ServeHTTP method
	ctrl.Ctx.Request = req
	ctrl.Ctx.ResponseWriter = rec

	// Call the method (don't pass req/rec directly)
	ctrl.GetImages()

	// Checking response
	assert.Equal(t, http.StatusOK, rec.Code)

	var response []models.CatImage
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "http://example.com/cat1.jpg", response[0].URL)

	// Assert that FetchCatImages was called once
	mockService.AssertExpectations(t)
}

// Test AddToFavorites
func TestAddToFavorites(t *testing.T) {
	mockService := new(MockBackendService)
	mockService.On("AddToFavorites", "1", "sub_123").Return(nil)

	ctrl := NewCatController(mockService)
	req := httptest.NewRequest(http.MethodPost, "/api/favourites", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Simulate the controller's ServeHTTP method
	ctrl.Ctx.Request = req
	ctrl.Ctx.ResponseWriter = rec

	ctrl.SaveFavorite()

	// Assert response
	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

// Test SaveVote
func TestSaveVote(t *testing.T) {
	mockService := new(MockBackendService)
	mockService.On("SaveVote", "1", 1, "sub_123").Return(nil)

	ctrl := NewCatController(mockService)
	req := httptest.NewRequest(http.MethodPost, "/api/vote/sub_123", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Simulate the controller's ServeHTTP method
	ctrl.Ctx.Request = req
	ctrl.Ctx.ResponseWriter = rec

	ctrl.SaveVote()

	// Assert response
	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

// Test error handling for AddToFavorites
func TestAddToFavoritesError(t *testing.T) {
	mockService := new(MockBackendService)
	mockService.On("AddToFavorites", "1", "sub_123").Return(errors.New("Failed to add to favorites"))

	ctrl := NewCatController(mockService)
	req := httptest.NewRequest(http.MethodPost, "/api/favourites", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Simulate the controller's ServeHTTP method
	ctrl.Ctx.Request = req
	ctrl.Ctx.ResponseWriter = rec

	ctrl.SaveFavorite()

	// Assert error response
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockService.AssertExpectations(t)
}

// Test error handling for SaveVote
func TestSaveVoteError(t *testing.T) {
	mockService := new(MockBackendService)
	mockService.On("SaveVote", "1", 1, "sub_123").Return(errors.New("Failed to save vote"))

	ctrl := NewCatController(mockService)
	req := httptest.NewRequest(http.MethodPost, "/api/vote/sub_123", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Simulate the controller's ServeHTTP method
	ctrl.Ctx.Request = req
	ctrl.Ctx.ResponseWriter = rec

	ctrl.SaveVote()

	// Assert error response
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockService.AssertExpectations(t)
}
