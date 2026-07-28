package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0xdiaz/gin-boilerplate/internal/app/controllers"
	"github.com/0xdiaz/gin-boilerplate/internal/app/dto"
	"github.com/0xdiaz/gin-boilerplate/internal/app/services/auth"
	"github.com/0xdiaz/gin-boilerplate/internal/domain/repositories"
	"github.com/0xdiaz/gin-boilerplate/tests/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupTestRouter creates a test Gin router
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

// setupAuthController creates an AuthController with real AuthService (for integration-style tests).
func setupAuthController() *controllers.AuthController {
	authService := auth.NewAuthService(repositories.NewUserRepository(), repositories.NewRefreshTokenRepository(), nil)
	return controllers.NewAuthController(authService)
}

// setupAuthControllerWithMock creates an AuthController with MockAuthServicer for unit tests.
func setupAuthControllerWithMock(mock *mocks.MockAuthServicer) *controllers.AuthController {
	return controllers.NewAuthController(mock)
}

// TestRegisterEndpoint tests the Register HTTP handler with mocked AuthServicer.
func TestRegisterEndpoint(t *testing.T) {
	mock := &mocks.MockAuthServicer{
		RegisterFunc: func(_ context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
			// Success for valid-looking request
			if req.Name != "" && req.Email != "" && len(req.Password) >= 8 {
				return &dto.AuthResponse{
					User:         dto.UserResponse{ID: 1, Name: req.Name, Email: req.Email},
					AccessToken:  "mock-token",
					RefreshToken: "mock-refresh",
					TokenType:    "Bearer",
				}, nil
			}
			return nil, auth.ErrEmailAlreadyExists
		},
	}
	ctrl := setupAuthControllerWithMock(mock)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Valid registration request",
			requestBody: dto.RegisterRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "SecurePass123!",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Missing required fields",
			requestBody: map[string]interface{}{
				"name": "John Doe",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, false, response["success"])
			},
		},
		{
			name: "Invalid email format",
			requestBody: dto.RegisterRequest{
				Name:     "John Doe",
				Email:    "invalid-email",
				Password: "SecurePass123!",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Password too short",
			requestBody: dto.RegisterRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "short",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid JSON",
			requestBody:    "invalid json string",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTestRouter()
			router.POST("/auth/register", ctrl.Register)

			// Create request body
			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code, "response status")
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

// TestLoginEndpoint tests the Login HTTP handler with mocked AuthService
func TestLoginEndpoint(t *testing.T) {
	mock := &mocks.MockAuthServicer{
		LoginFunc: func(_ context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
			if req.Email != "" && req.Password != "" {
				return &dto.AuthResponse{
					User:         dto.UserResponse{ID: 1, Name: "Test", Email: req.Email},
					AccessToken:  "mock-token",
					RefreshToken: "mock-refresh",
					TokenType:    "Bearer",
				}, nil
			}
			return nil, auth.ErrInvalidCredentials
		},
	}
	ctrl := setupAuthControllerWithMock(mock)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Valid login request format",
			requestBody: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Missing email",
			requestBody: map[string]interface{}{
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Missing password",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Empty request body",
			requestBody:    map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid JSON",
			requestBody:    "not a json",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTestRouter()
			router.POST("/auth/login", ctrl.Login)

			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code, "response status")
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

// TestRegisterResponseStructure tests the response structure with mocked AuthService
func TestRegisterResponseStructure(t *testing.T) {
	mock := &mocks.MockAuthServicer{
		RegisterFunc: func(_ context.Context, _ *dto.RegisterRequest) (*dto.AuthResponse, error) {
			return &dto.AuthResponse{
				User:         dto.UserResponse{ID: 1, Name: "Test User", Email: "testuser@example.com"},
				AccessToken:  "mock-access-token",
				RefreshToken: "mock-refresh-token",
				TokenType:    "Bearer",
			}, nil
		},
	}
	ctrl := setupAuthControllerWithMock(mock)

	router := setupTestRouter()
	router.POST("/auth/register", ctrl.Register)

	requestBody := dto.RegisterRequest{
		Name:     "Test User",
		Email:    "testuser@example.com",
		Password: "SecurePass123!",
	}

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code == http.StatusCreated {
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check new standard response structure
		assert.Contains(t, response, "success")
		assert.Contains(t, response, "message")
		assert.Contains(t, response, "data")
		assert.Equal(t, true, response["success"])

		data := response["data"].(map[string]interface{})
		assert.Contains(t, data, "user")
		assert.Contains(t, data, "access_token")
		assert.Contains(t, data, "token_type")
	}
}

// TestProfileEndpoint tests the Profile HTTP handler with mocked AuthServicer.
func TestProfileEndpoint(t *testing.T) {
	tests := []struct {
		name           string
		mock           *mocks.MockAuthServicer
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Existing user returns real profile data",
			mock: &mocks.MockAuthServicer{
				GetProfileFunc: func(_ context.Context, userID uint) (*dto.UserResponse, error) {
					return &dto.UserResponse{ID: userID, Name: "Test User", Email: "test@example.com"}, nil
				},
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				data := response["data"].(map[string]interface{})
				assert.Equal(t, "Test User", data["name"])
				assert.Equal(t, "test@example.com", data["email"])
			},
		},
		{
			name: "User not found returns 404",
			mock: &mocks.MockAuthServicer{
				GetProfileFunc: func(_ context.Context, _ uint) (*dto.UserResponse, error) {
					return nil, auth.ErrUserNotFound
				},
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := setupAuthControllerWithMock(tt.mock)
			router := setupTestRouter()
			router.GET("/profile", func(c *gin.Context) {
				c.Set("user_id", uint(1))
				ctrl.Profile(c)
			})

			req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code, "response status")
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

// Example: Benchmark test for controller
func BenchmarkRegisterEndpoint(b *testing.B) {
	router := setupTestRouter()
	authController := setupAuthController()
	router.POST("/auth/register", authController.Register)

	requestBody := dto.RegisterRequest{
		Name:     "Benchmark User",
		Email:    "bench@example.com",
		Password: "SecurePass123!",
	}
	body, _ := json.Marshal(requestBody)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
