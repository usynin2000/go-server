package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware логирует все HTTP запросы
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Вызываем следующий обработчик
		next.ServeHTTP(w, r)

		// Логируем запрос
		log.Printf(
			"%s %s %s %v",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			time.Since(start),
		)
	})
}
