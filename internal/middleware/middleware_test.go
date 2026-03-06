package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// TestSecurityHeaders тестирует заголовки безопасности.
func TestSecurityHeaders(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		header string
		value  string
	}{
		{"X-Content-Type-Options", "X-Content-Type-Options", "nosniff"},
		{"X-Frame-Options", "X-Frame-Options", "DENY"},
		{"X-XSS-Protection", "X-XSS-Protection", "1; mode=block"},
		{"Referrer-Policy", "Referrer-Policy", "strict-origin-when-cross-origin"},
		{"Content-Security-Policy", "Content-Security-Policy", "default-src 'self'"},
		{"Cache-Control", "Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate"},
		{"Pragma", "Pragma", "no-cache"},
		{"Expires", "Expires", "0"},
	}

	handler := SecurityHeaders()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if got := w.Header().Get(tt.header); got != tt.value {
				t.Errorf("expected %s header %q, got %q", tt.header, tt.value, got)
			}
		})
	}
}

// TestContentTypeValidator тестирует валидацию Content-Type.
func TestContentTypeValidator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		method      string
		contentType string
		wantStatus  int
	}{
		// POST запросы
		{
			name:        "POST with application/json",
			method:      http.MethodPost,
			contentType: "application/json",
			wantStatus:  http.StatusOK,
		},
		{
			name:        "POST with application/json; charset=utf-8",
			method:      http.MethodPost,
			contentType: "application/json; charset=utf-8",
			wantStatus:  http.StatusOK,
		},
		{
			name:        "POST with text/plain",
			method:      http.MethodPost,
			contentType: "text/plain",
			wantStatus:  http.StatusUnsupportedMediaType,
		},
		{
			name:        "POST without Content-Type",
			method:      http.MethodPost,
			contentType: "",
			wantStatus:  http.StatusUnsupportedMediaType,
		},
		// GET запросы (не требуют Content-Type)
		{
			name:        "GET without Content-Type",
			method:      http.MethodGet,
			contentType: "",
			wantStatus:  http.StatusOK,
		},
		{
			name:        "GET with application/json",
			method:      http.MethodGet,
			contentType: "application/json",
			wantStatus:  http.StatusOK,
		},
		// OPTIONS запросы (не требуют Content-Type)
		{
			name:        "OPTIONS without Content-Type",
			method:      http.MethodOptions,
			contentType: "",
			wantStatus:  http.StatusOK,
		},
		// PUT запросы (требуют Content-Type)
		{
			name:        "PUT with application/json",
			method:      http.MethodPut,
			contentType: "application/json",
			wantStatus:  http.StatusOK,
		},
		{
			name:        "PUT without Content-Type",
			method:      http.MethodPut,
			contentType: "",
			wantStatus:  http.StatusUnsupportedMediaType,
		},
		// DELETE запросы (не требуют Content-Type)
		{
			name:        "DELETE without Content-Type",
			method:      http.MethodDelete,
			contentType: "",
			wantStatus:  http.StatusOK,
		},
		{
			name:        "DELETE with application/json",
			method:      http.MethodDelete,
			contentType: "application/json",
			wantStatus:  http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var handlerCalled bool
			handler := ContentTypeValidator("application/json")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(tt.method, "/", nil)
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
			
			// Проверяем что handler был вызван при успешном статусе
			if tt.wantStatus == http.StatusOK && !handlerCalled {
				t.Error("expected handler to be called")
			}
		})
	}
}

// TestRequestLimit тестирует ограничение размера запроса.
func TestRequestLimit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		body       string
		maxBytes   int64
		wantStatus int
		readError  bool // Ожидается ошибка при чтении
	}{
		{
			name:       "body under limit",
			body:       `{"name": "test"}`,
			maxBytes:   1024,
			wantStatus: http.StatusOK,
			readError:  false,
		},
		{
			name:       "body over limit",
			body:       strings.Repeat("a", 200),
			maxBytes:   100,
			wantStatus: http.StatusOK,
			readError:  true, // Ошибка возникнет при чтении тела
		},
		{
			name:       "empty body",
			body:       "",
			maxBytes:   1024,
			wantStatus: http.StatusOK,
			readError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var handlerCalled bool
			handler := RequestLimit(tt.maxBytes)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
				// Пытаемся прочитать тело - здесь возникнет ошибка если превышен лимит
				_, err := r.Body.Read(make([]byte, tt.maxBytes+1))
				if tt.readError && err == nil {
					t.Error("expected error when reading body over limit")
				}
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
			if !handlerCalled {
				t.Error("expected handler to be called")
			}
		})
	}
}

// TestRequestLogger тестирует логирование запросов.
func TestRequestLogger(t *testing.T) {
	t.Parallel()

	logger, _ := zap.NewDevelopment()

	handler := RequestLogger(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

// TestChain тестирует цепочку middleware.
func TestChain(t *testing.T) {
	t.Parallel()

	logger, _ := zap.NewDevelopment()

	var handlerCalled bool
	handler := Chain(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if !handlerCalled {
		t.Error("expected handler to be called")
	}
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Проверяем заголовки безопасности
	headers := []string{
		"X-Content-Type-Options",
		"X-Frame-Options",
		"X-XSS-Protection",
	}
	for _, header := range headers {
		if w.Header().Get(header) == "" {
			t.Errorf("expected %s header to be set", header)
		}
	}
}

// TestRequestLogger_PanicRecovery тестирует восстановление после паники.
func TestRequestLogger_PanicRecovery(t *testing.T) {
	t.Parallel()

	logger, _ := zap.NewDevelopment()

	handler := RequestLogger(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Recoverer должен перехватить панику
	fullHandler := chiMiddleware.Recoverer(handler)

	defer func() {
		if r := recover(); r != nil {
			t.Error("panic was not recovered")
		}
	}()

	fullHandler.ServeHTTP(w, req)
}

// BenchmarkSecurityHeaders бенчмарк middleware заголовков.
func BenchmarkSecurityHeaders(b *testing.B) {
	handler := SecurityHeaders()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(w, req)
	}
}

// BenchmarkContentTypeValidator бенчмарк валидации Content-Type.
func BenchmarkContentTypeValidator(b *testing.B) {
	handler := ContentTypeValidator("application/json")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(w, req)
	}
}
