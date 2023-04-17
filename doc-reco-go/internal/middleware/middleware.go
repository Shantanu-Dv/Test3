package middleware

import (
	sentryHttp "github.com/getsentry/sentry-go/http"
	"log"
	"net/http"
	"time"
)

func SetHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func SentryMiddleware(next http.Handler) http.Handler {
	sentryHandler := sentryHttp.New(sentryHttp.Options{Repanic: true})
	return sentryHandler.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s: %s took %d ms", r.RemoteAddr, r.RequestURI, time.Since(t).Milliseconds())
	})
}
