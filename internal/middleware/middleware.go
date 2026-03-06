// Package middleware содержит HTTP middleware для обработки запросов.
package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// Chain возвращает цепочку стандартных middleware.
func Chain(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return middleware.Recoverer(
			RequestLogger(logger)(
				SecurityHeaders()(next),
			),
		)
	}
}

// RequestLogger создаёт middleware для логирования запросов с zap.
func RequestLogger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				logger.Info("request completed",
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.String("remote", r.RemoteAddr),
					zap.Int("status", ww.Status()),
					zap.Duration("duration", time.Since(start)),
					zap.Int("bytes", ww.BytesWritten()),
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}

// SecurityHeaders создаёт middleware для заголовков безопасности.
func SecurityHeaders() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			w.Header().Set("Content-Security-Policy", "default-src 'self'")
			w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
			next.ServeHTTP(w, r)
		})
	}
}

// ContentTypeValidator создаёт middleware для валидации Content-Type.
func ContentTypeValidator(validTypes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Пропускаем методы без тела запроса
			if r.Method == http.MethodGet || r.Method == http.MethodHead || 
				r.Method == http.MethodOptions || r.Method == http.MethodDelete {
				next.ServeHTTP(w, r)
				return
			}

			contentType := r.Header.Get("Content-Type")
			for _, valid := range validTypes {
				if contentType == valid {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Разрешаем application/json с charset
			if contentType != "" && len(contentType) >= 16 && contentType[:16] == "application/json" {
				next.ServeHTTP(w, r)
				return
			}

			http.Error(w, `{"error":"Unsupported Media Type. Expected application/json"}`, http.StatusUnsupportedMediaType)
		})
	}
}

// RequestLimit создаёт middleware для ограничения размера тела запроса.
func RequestLimit(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next.ServeHTTP(w, r)
		})
	}
}
