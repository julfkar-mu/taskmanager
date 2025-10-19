package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"taskmanager/constants"
	"taskmanager/errors"
	"taskmanager/models"
	"taskmanager/testutils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTaskService is a mock implementation of TaskService for testing
type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) GetTasks() []models.Task {
	args := m.Called()
	return args.Get(0).([]models.Task)
}

func (m *MockTaskService) GetTask(id string) (models.Task, error) {
	args := m.Called(id)
	return args.Get(0).(models.Task), args.Error(1)
}

func (m *MockTaskService) CreateTask(task models.Task) (models.Task, error) {
	args := m.Called(task)
	return args.Get(0).(models.Task), args.Error(1)
}

func (m *MockTaskService) UpdateTask(id string, task models.Task) (models.Task, error) {
	args := m.Called(id, task)
	return args.Get(0).(models.Task), args.Error(1)
}

func (m *MockTaskService) DeleteTask(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestGetTasks(t *testing.T) {
	mockService := new(MockTaskService)
	Setup(mockService)

	tests := []struct {
		name           string
		mockTasks      []models.Task
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "Empty tasks list",
			mockTasks:      []models.Task{},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name: "Tasks list with data",
			mockTasks: []models.Task{
				testutils.CreateTestTask(),
				testutils.CreateTestTask(),
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.On("GetTasks").Return(tt.mockTasks)

			router := setupTestRouter()
			router.GET("/tasks", GetTasks)

			req, _ := http.NewRequest("GET", "/tasks", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			data := response["data"].([]interface{})
			assert.Equal(t, tt.expectedCount, len(data))
			assert.Equal(t, float64(tt.expectedCount), response["count"])

			mockService.AssertExpectations(t)
		})
	}
}

func TestGetTaskByID(t *testing.T) {
	mockService := new(MockTaskService)
	Setup(mockService)

	tests := []struct {
		name           string
		taskID         string
		mockTask       models.Task
		mockError      error
		expectedStatus int
	}{
		{
			name:           "Valid task ID",
			taskID:         "test-id",
			mockTask:       testutils.CreateTestTask(),
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Non-existent task ID",
			taskID:         "non-existent",
			mockTask:       models.Task{},
			mockError:      errors.NewNotFoundError("Task"),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.On("GetTask", tt.taskID).Return(tt.mockTask, tt.mockError)

			router := setupTestRouter()
			router.GET("/tasks/:id", GetTaskByID)

			req, _ := http.NewRequest("GET", "/tasks/"+tt.taskID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "data")
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestCreateTask(t *testing.T) {
	mockService := new(MockTaskService)
	Setup(mockService)

	tests := []struct {
		name           string
		requestBody    models.Task
		mockTask       models.Task
		mockError      error
		expectedStatus int
	}{
		{
			name:        "Valid task creation",
			requestBody: testutils.CreateTestTask(),
			mockTask:    testutils.CreateTestTask(),
			mockError:   nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:        "Invalid task data",
			requestBody: testutils.CreateInvalidTask(),
			mockTask:    models.Task{},
			mockError:   errors.NewValidationError("title", constants.ValidationTitleRequired),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockError == nil {
				mockService.On("CreateTask", mock.AnythingOfType("models.Task")).Return(tt.mockTask, tt.mockError)
			} else {
				mockService.On("CreateTask", mock.AnythingOfType("models.Task")).Return(tt.mockTask, tt.mockError)
			}

			router := setupTestRouter()
			router.POST("/tasks", CreateTask)

			jsonBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "data")
				assert.Contains(t, response, "message")
				assert.Equal(t, constants.MessageTaskCreated, response["message"])
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestUpdateTask(t *testing.T) {
	mockService := new(MockTaskService)
	Setup(mockService)

	tests := []struct {
		name           string
		taskID         string
		requestBody    models.Task
		mockTask       models.Task
		mockError      error
		expectedStatus int
	}{
		{
			name:        "Valid task update",
			taskID:      "test-id",
			requestBody: testutils.CreateTestTask(),
			mockTask:    testutils.CreateTestTask(),
			mockError:   nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:        "Non-existent task",
			taskID:      "non-existent",
			requestBody: testutils.CreateTestTask(),
			mockTask:    models.Task{},
			mockError:   errors.NewNotFoundError("Task"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:        "Invalid task data",
			taskID:      "test-id",
			requestBody: testutils.CreateInvalidTask(),
			mockTask:    models.Task{},
			mockError:   errors.NewValidationError("title", constants.ValidationTitleRequired),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.On("UpdateTask", tt.taskID, mock.AnythingOfType("models.Task")).Return(tt.mockTask, tt.mockError)

			router := setupTestRouter()
			router.PUT("/tasks/:id", UpdateTask)

			jsonBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("PUT", "/tasks/"+tt.taskID, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "data")
				assert.Contains(t, response, "message")
				assert.Equal(t, constants.MessageTaskUpdated, response["message"])
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestDeleteTask(t *testing.T) {
	mockService := new(MockTaskService)
	Setup(mockService)

	tests := []struct {
		name           string
		taskID         string
		mockError      error
		expectedStatus int
	}{
		{
			name:           "Valid task deletion",
			taskID:         "test-id",
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Non-existent task",
			taskID:         "non-existent",
			mockError:      errors.NewNotFoundError("Task"),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.On("DeleteTask", tt.taskID).Return(tt.mockError)

			router := setupTestRouter()
			router.DELETE("/tasks/:id", DeleteTask)

			req, _ := http.NewRequest("DELETE", "/tasks/"+tt.taskID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "message")
				assert.Equal(t, constants.MessageTaskDeleted, response["message"])
			}

			mockService.AssertExpectations(t)
		})
	}
}
