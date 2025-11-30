package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"repo-guardian/internal/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) Register(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserUsecase) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	args := m.Called(ctx, id)
	user, ok := args.Get(0).(*domain.User)
	if !ok {
		return nil, args.Error(1)
	}
	return user, args.Error(1)
}

func (m *MockUserUsecase) DeleteUser(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestNewUserHandler(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(MockUserUsecase)

	NewUserHandler(app, mockUsecase)

	routes := app.GetRoutes()
	var postUsers, getUsers, deleteUsers bool
	for _, r := range routes {
		switch {
		case r.Method == http.MethodPost && r.Path == "/users":
			postUsers = true
		case r.Method == http.MethodGet && strings.HasPrefix(r.Path, "/users/:id"):
			getUsers = true
		case r.Method == http.MethodDelete && strings.HasPrefix(r.Path, "/users/:id"):
			deleteUsers = true
		}
	}

	assert.True(t, postUsers, "POST /users route not registered")
	assert.True(t, getUsers, "GET /users/:id route not registered")
	assert.True(t, deleteUsers, "DELETE /users/:id route not registered")
}

func TestUserHandler_Register(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(MockUserUsecase)
	handler := &UserHandler{UserUsecase: mockUsecase}
	app.Post("/users", handler.Register)

	t.Run("Success", func(t *testing.T) {
		user := domain.User{Name: "Test User", Email: "test@example.com"}
		userJSON, _ := json.Marshal(user)

		mockUsecase.On("Register", mock.Anything, &user).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(userJSON))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var registeredUser domain.User
		err = json.Unmarshal(body, &registeredUser)
		assert.NoError(t, err)

		assert.Equal(t, user.Name, registeredUser.Name)
		assert.Equal(t, user.Email, registeredUser.Email)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("BadRequest_InvalidBody", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InternalServerError_UsecaseError", func(t *testing.T) {
		user := domain.User{Name: "Test User", Email: "test@example.com"}
		userJSON, _ := json.Marshal(user)

		mockUsecase.On("Register", mock.Anything, &user).Return(errors.New("usecase error")).Once()

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(userJSON))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockUsecase.AssertExpectations(t)
	})
}

func TestUserHandler_GetUser(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(MockUserUsecase)
	handler := &UserHandler{UserUsecase: mockUsecase}
	app.Get("/users/:id", handler.GetUser)

	t.Run("Success", func(t *testing.T) {
		id := int64(1)
		user := &domain.User{ID: id, Name: "Test User", Email: "test@example.com"}
		mockUsecase.On("GetUser", mock.Anything, id).Return(user, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/users/1", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var retrievedUser domain.User
		err = json.Unmarshal(body, &retrievedUser)
		assert.NoError(t, err)

		assert.Equal(t, user.ID, retrievedUser.ID)
		assert.Equal(t, user.Name, retrievedUser.Name)
		assert.Equal(t, user.Email, retrievedUser.Email)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("BadRequest_InvalidID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/invalid", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "Invalid ID")
	})

	t.Run("NotFound_UsecaseError", func(t *testing.T) {
		id := int64(1)
		mockUsecase.On("GetUser", mock.Anything, id).Return(nil, errors.New("user not found")).Once()

		req := httptest.NewRequest(http.MethodGet, "/users/1", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "user not found")
		mockUsecase.AssertExpectations(t)
	})
}

func TestUserHandler_DeleteUser(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(MockUserUsecase)
	handler := &UserHandler{UserUsecase: mockUsecase}
	app.Delete("/users/:id", handler.DeleteUser)

	t.Run("Success", func(t *testing.T) {
		id := int64(1)
		mockUsecase.On("DeleteUser", mock.Anything, id).Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("BadRequest_InvalidID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/users/invalid", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "Invalid ID")
	})

	t.Run("InternalServerError_UsecaseError", func(t *testing.T) {
		id := int64(1)
		mockUsecase.On("DeleteUser", mock.Anything, id).Return(errors.New("delete failed")).Once()

		req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "delete failed")
		mockUsecase.AssertExpectations(t)
	})
}

func TestUserHandler_Register_Integration(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(MockUserUsecase)
	handler := &UserHandler{UserUsecase: mockUsecase}
	app.Post("/users", handler.Register)

	t.Run("Register User - Valid Input", func(t *testing.T) {
		// Define a valid user
		user := domain.User{Name: "John Doe", Email: "john.doe@example.com"}
		userJSON, err := json.Marshal(user)
		assert.NoError(t, err)

		// Mock the Usecase to return nil (no error)
		mockUsecase.On("Register", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
			return u.Name == user.Name && u.Email == user.Email
		})).Return(nil).Once()

		// Create a request to the endpoint
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(userJSON))
		req.Header.Set("Content-Type", "application/json")

		// Perform the request
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		// Unmarshal the response into a User struct
		var responseUser domain.User
		err = json.Unmarshal(body, &responseUser)
		assert.NoError(t, err)

		// Assert the response matches the expected user
		assert.Equal(t, user.Name, responseUser.Name)
		assert.Equal(t, user.Email, responseUser.Email)

		mockUsecase.AssertExpectations(t)
	})
}

func TestUserHandler_GetUser_Integration(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(MockUserUsecase)
	handler := &UserHandler{UserUsecase: mockUsecase}
	app.Get("/users/:id", handler.GetUser)

	t.Run("Get User - Valid ID", func(t *testing.T) {
		// Define a valid user ID
		userID := int64(123)

		// Define the user that the Usecase will return
		expectedUser := &domain.User{ID: userID, Name: "Jane Doe", Email: "jane.doe@example.com"}

		// Mock the Usecase to return the expected user
		mockUsecase.On("GetUser", mock.Anything, userID).Return(expectedUser, nil).Once()

		// Create a request to the endpoint
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%d", userID), nil)

		// Perform the request
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		// Unmarshal the response into a User struct
		var responseUser domain.User
		err = json.Unmarshal(body, &responseUser)
		assert.NoError(t, err)

		// Assert the response matches the expected user
		assert.Equal(t, expectedUser.ID, responseUser.ID)
		assert.Equal(t, expectedUser.Name, responseUser.Name)
		assert.Equal(t, expectedUser.Email, responseUser.Email)

		mockUsecase.AssertExpectations(t)
	})

	t.Run("Get User - Invalid ID (String)", func(t *testing.T) {
		// Create a request with an invalid ID
		req := httptest.NewRequest(http.MethodGet, "/users/invalid_id", nil)

		// Perform the request
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		// Assert the response body contains the error message
		assert.Contains(t, string(body), "Invalid ID")

		mockUsecase.AssertExpectations(t) //Expect no calls to the mock usecase
	})

	t.Run("Get User - Not Found", func(t *testing.T) {
		// Define a user ID that will return a "not found" error
		userID := int64(999)

		// Mock the Usecase to return an error
		mockUsecase.On("GetUser", mock.Anything, userID).Return(nil, errors.New("user not found")).Once()

		// Create a request to the endpoint
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%d", userID), nil)

		// Perform the request
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		// Assert the response body contains the error message
		assert.Contains(t, string(body), "user not found")

		mockUsecase.AssertExpectations(t)
	})
}

func TestUserHandler_DeleteUser_Integration(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(MockUserUsecase)
	handler := &UserHandler{UserUsecase: mockUsecase}
	app.Delete("/users/:id", handler.DeleteUser)

	t.Run("Delete User - Valid ID", func(t *testing.T) {
		// Define a valid user ID
		userID := int64(456)

		// Mock the Usecase to return nil (no error)
		mockUsecase.On("DeleteUser", mock.Anything, userID).Return(nil).Once()

		// Create a request to the endpoint
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%d", userID), nil)

		// Perform the request
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		mockUsecase.AssertExpectations(t)
	})

	t.Run("Delete User - Invalid ID (String)", func(t *testing.T) {
		// Create a request with an invalid ID
		req := httptest.NewRequest(http.MethodDelete, "/users/invalid_id", nil)

		// Perform the request
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		// Assert the response body contains the error message
		assert.Contains(t, string(body), "Invalid ID")

		mockUsecase.AssertExpectations(t)
	})

	t.Run("Delete User - Internal Server Error", func(t *testing.T) {
		// Define a user ID that will return an error
		userID := int64(789)

		// Mock the Usecase to return an error
		mockUsecase.On("DeleteUser", mock.Anything, userID).Return(errors.New("delete failed")).Once()

		// Create a request to the endpoint
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%d", userID), nil)

		// Perform the request
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		// Assert the response body contains the error message
		assert.Contains(t, string(body), "delete failed")

		mockUsecase.AssertExpectations(t)
	})
}

func TestUserHandler_DeleteUser_Concurrent(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(MockUserUsecase)
	handler := &UserHandler{UserUsecase: mockUsecase}
	app.Delete("/users/:id", handler.DeleteUser)

	numRequests := 10
	userID := int64(123)

	// Mock the Usecase to return nil (no error) for all calls
	mockUsecase.On("DeleteUser", mock.Anything, userID).Return(nil).Times(numRequests)

	errChan := make(chan error, numRequests)

	// Create a worker pool to handle the requests concurrently
	for i := 0; i < numRequests; i++ {
		go func(requestNum int) {
			// Create a request to the endpoint
			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%d", userID), nil)

			// Perform the request
			resp, err := app.Test(req)
			if err != nil {
				errChan <- fmt.Errorf("request %d failed: %w", requestNum, err)
				return
			}
			if resp.StatusCode != http.StatusNoContent {
				errChan <- fmt.Errorf("request %d failed: unexpected status code %d", requestNum, resp.StatusCode)
				return
			}
			errChan <- nil // No error
		}(i)
	}

	// Wait for all requests to complete and check for errors
	for i := 0; i < numRequests; i++ {
		err := <-errChan
		if err != nil {
			t.Error(err)
		}
	}

	mockUsecase.AssertExpectations(t)
	close(errChan)
}

func TestGetUserHandler_InvalidID(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(MockUserUsecase)
	handler := &UserHandler{UserUsecase: mockUsecase}
	app.Get("/users/:id", handler.GetUser)

	req := httptest.NewRequest(http.MethodGet, "/users/abc", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestDeleteUserHandler_InvalidID(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(MockUserUsecase)
	handler := &UserHandler{UserUsecase: mockUsecase}
	app.Delete("/users/:id", handler.DeleteUser)

	req := httptest.NewRequest(http.MethodDelete, "/users/abc", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestParseIntError(t *testing.T) {
	// Test Case 1: Empty string
	_, err := strconv.ParseInt("", 10, 64)
	assert.Error(t, err)

	// Test Case 2: Invalid character
	_, err = strconv.ParseInt("12a3", 10, 64)
	assert.Error(t, err)

	// Test Case 3: Number too large
	_, err = strconv.ParseInt("9223372036854775808", 10, 64) // One more than max int64
	assert.Error(t, err)

	// Test Case 4: Negative number (not applicable for uint64, but tests the function)
	_, err = strconv.ParseInt("-1", 10, 64)
	assert.Error(t, err)

	// Test Case 5: Valid number
	num, err := strconv.ParseInt("12345", 10, 64)
	assert.NoError(t, err)
	assert.Equal(t, int64(12345), num)
}
