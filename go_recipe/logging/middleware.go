package logging

import (
	"net/http"

	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
)

type key int

const (
	requestIDKey string = "X-Request-ID"
	userAgentKey string = "User-Agent"
)

// WithRequestID returns a log entry with the request id.
func WithRequestID(key string, r *http.Request) *logrus.Entry {
	return logrus.WithField(key, r.Header.Get(requestIDKey))
}

func getOrAssignRequestID(r *http.Request) string {
	id := r.Header.Get(requestIDKey)
	if id == "" {
		id = xid.New().String()
		r.Header.Set(requestIDKey, id)
	}
	return id
}

type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLogResponseWriter(w http.ResponseWriter) *logResponseWriter {
	return &logResponseWriter{w, -1}
}

// WriteHeader records response code for final logging.
func (lrw *logResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *logResponseWriter) logResponse(r *http.Request) {
	WithRequestID("rid", r).WithField("status", lrw.statusCode).Infof("%s %s", r.Method, r.URL.Path)
}

// Middleware logs the request and its response. Also sets request ID into context for logging.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log incoming request with serial ID.
		reqID := getOrAssignRequestID(r)
		logrus.WithFields(logrus.Fields{
			"rid":   reqID,
			"addr":  r.RemoteAddr,
			"agent": r.Header.Get(userAgentKey),
		}).Infof("%s %s", r.Method, r.URL.Path)

		// Log response with status code on completion.
		lrw := newLogResponseWriter(w)
		defer lrw.logResponse(r)

		next.ServeHTTP(lrw, r)
	})
}
