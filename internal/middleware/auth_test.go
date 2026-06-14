package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hawk-roy/Night-Hawk/internal/auth"
	"github.com/hawk-roy/Night-Hawk/internal/response"
)

type testBody struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func decodeTestBody(t *testing.T, w *httptest.ResponseRecorder) testBody {
	t.Helper()

	var body testBody
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	return body
}

func newProtectedRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/protected", AuthMiddleware(), func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")

		response.Success(c, gin.H{
			"user_id":  userID,
			"username": username,
		})
	})

	return r
}

func TestAuthMiddlewareMissingToken(t *testing.T) {
	r := newProtectedRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected http status: got %d want %d", w.Code, http.StatusUnauthorized)
	}

	body := decodeTestBody(t, w)
	if body.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected code: got %d want %d", body.Code, http.StatusUnauthorized)
	}
	if body.Message != "unauthorized" {
		t.Fatalf("unexpected message: got %q want %q", body.Message, "unauthorized")
	}
	if string(body.Data) != "null" {
		t.Fatalf("expected data to be null, got %s", string(body.Data))
	}
}

func TestAuthMiddlewareInvalidToken(t *testing.T) {
	r := newProtectedRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected http status: got %d want %d", w.Code, http.StatusUnauthorized)
	}

	body := decodeTestBody(t, w)
	if body.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected code: got %d want %d", body.Code, http.StatusUnauthorized)
	}
	if body.Message != "unauthorized" {
		t.Fatalf("unexpected message: got %q want %q", body.Message, "unauthorized")
	}
	if string(body.Data) != "null" {
		t.Fatalf("expected data to be null, got %s", string(body.Data))
	}
}

func TestAuthMiddlewareValidToken(t *testing.T) {
	r := newProtectedRouter()

	token, err := auth.GenerateToken(123, "jwt_user")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected http status: got %d want %d", w.Code, http.StatusOK)
	}

	body := decodeTestBody(t, w)
	if body.Code != 0 {
		t.Fatalf("unexpected code: got %d want %d", body.Code, 0)
	}
	if body.Message != "success" {
		t.Fatalf("unexpected message: got %q want %q", body.Message, "success")
	}

	var data map[string]any
	if err := json.Unmarshal(body.Data, &data); err != nil {
		t.Fatalf("failed to decode data: %v", err)
	}

	if got, ok := data["user_id"].(float64); !ok || got != 123 {
		t.Fatalf("unexpected user_id: %+v", data["user_id"])
	}
	if got, ok := data["username"].(string); !ok || got != "jwt_user" {
		t.Fatalf("unexpected username: %+v", data["username"])
	}
}
