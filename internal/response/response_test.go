package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type body struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func decodeBody(t *testing.T, w *httptest.ResponseRecorder) body {
	t.Helper()

	var resp body
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	return resp
}

func TestSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Success(c, gin.H{"ok": true})

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected http status: got %d want %d", w.Code, http.StatusOK)
	}

	resp := decodeBody(t, w)

	if resp.Code != 0 {
		t.Fatalf("unexpected code: got %d want %d", resp.Code, 0)
	}
	if resp.Message != "success" {
		t.Fatalf("unexpected message: got %q want %q", resp.Message, "success")
	}
	if len(resp.Data) == 0 || string(resp.Data) == "null" {
		t.Fatalf("expected non-empty data, got %s", string(resp.Data))
	}

	var data map[string]any
	if err := json.Unmarshal(resp.Data, &data); err != nil {
		t.Fatalf("failed to decode data: %v", err)
	}
	if ok, exists := data["ok"].(bool); !exists || !ok {
		t.Fatalf("unexpected data: %+v", data)
	}
}

func TestSuccessWithStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SuccessWithStatus(c, http.StatusCreated, gin.H{"id": 1})

	if w.Code != http.StatusCreated {
		t.Fatalf("unexpected http status: got %d want %d", w.Code, http.StatusCreated)
	}

	resp := decodeBody(t, w)

	if resp.Code != 0 {
		t.Fatalf("unexpected code: got %d want %d", resp.Code, 0)
	}
	if resp.Message != "success" {
		t.Fatalf("unexpected message: got %q want %q", resp.Message, "success")
	}

	var data map[string]any
	if err := json.Unmarshal(resp.Data, &data); err != nil {
		t.Fatalf("failed to decode data: %v", err)
	}
	if id, exists := data["id"].(float64); !exists || id != 1 {
		t.Fatalf("unexpected data: %+v", data)
	}
}

func TestError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Error(c, http.StatusBadRequest, http.StatusBadRequest, "invalid request")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("unexpected http status: got %d want %d", w.Code, http.StatusBadRequest)
	}

	resp := decodeBody(t, w)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("unexpected code: got %d want %d", resp.Code, http.StatusBadRequest)
	}
	if resp.Message != "invalid request" {
		t.Fatalf("unexpected message: got %q want %q", resp.Message, "invalid request")
	}
	if string(resp.Data) != "null" {
		t.Fatalf("expected data to be null, got %s", string(resp.Data))
	}
}
