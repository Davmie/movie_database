package middleware

import (
	"fmt"
	"net/http"

	"intern/pkg/logger"
)

func Panic(logger logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				logger.Errorw("Server paniced",
					"method", r.Method,
					"remote_addr", r.RemoteAddr,
					"url", r.URL.Path,
					"error", err,
				)

				fmt.Println("recovered", err)
				http.Error(w, "Internal server error", 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
