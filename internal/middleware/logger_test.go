package middleware

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hawk-roy/Night-Hawk/internal/response"
)

type loggerBody struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func decodeLoggerBody(t *testing.T, w *httptest.ResponseRecorder) loggerBody {
	t.Helper()

	var body loggerBody
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	return body
}

func newLoggerRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(RequestLogger())
	r.GET("/ping", func(c *gin.Context) {
		response.Success(c, gin.H{"pong": true})
	})
	return r
}

func TestRequestLoggerGeneratesRequestID(t *testing.T) {
	r := newLoggerRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected http status: got %d want %d", w.Code, http.StatusOK)
	}

	if got := w.Header().Get("X-Request-ID"); got == "" {
		t.Fatal("expected X-Request-ID header to be set")
	}

	body := decodeLoggerBody(t, w)
	if body.Code != 0 {
		t.Fatalf("unexpected code: got %d want %d", body.Code, 0)
	}
}

func TestRequestLoggerPassesThroughRequestID(t *testing.T) {
	r := newLoggerRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("X-Request-ID", "test-request-id")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected http status: got %d want %d", w.Code, http.StatusOK)
	}

	if got := w.Header().Get("X-Request-ID"); got != "test-request-id" {
		t.Fatalf("unexpected X-Request-ID: got %q want %q", got, "test-request-id")
	}
}

func TestRequestLoggerWritesLog(t *testing.T) {
	r := newLoggerRouter()

	var buf bytes.Buffer
	oldOutput := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(oldOutput)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("X-Request-ID", "test-request-id")
	r.ServeHTTP(w, req)

	logLine := buf.String()
	if !strings.Contains(logLine, "request_id=test-request-id") {
		t.Fatalf("expected log to contain request_id, got %q", logLine)
	}
	if !strings.Contains(logLine, "method=GET") {
		t.Fatalf("expected log to contain method, got %q", logLine)
	}
	if !strings.Contains(logLine, "path=/ping") {
		t.Fatalf("expected log to contain path, got %q", logLine)
	}
	if !strings.Contains(logLine, "status=200") {
		t.Fatalf("expected log to contain status, got %q", logLine)
	}
}
