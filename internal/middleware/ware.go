package middleware

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func ScanTrafic(log *logrus.Logger) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			log.WithFields(logrus.Fields{
				"method": r.Method,
				"url":    r.URL,
				"addr":   r.RemoteAddr,
				"time":   time.Now().Format(time.RFC3339),
			}).Info("Request received")

			h.ServeHTTP(w, r)
		})
	}
}
