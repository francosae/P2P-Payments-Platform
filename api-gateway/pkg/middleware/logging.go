package middleware

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/docs" || c.Request.URL.Path == "/docs/swagger.json" || c.Request.URL.Path == "/swagger/index.html" {
			c.Next()
			return
		}

		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		w := &responseWriter{body: new(bytes.Buffer), ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		if status >= http.StatusBadRequest {
			log.Error().
				Str("method", method).
				Str("path", path).
				Int("status", status).
				Str("response", w.body.String()).
				Dur("duration", duration).
				Msg("Failed request details")
		} else {
			log.Info().
				Str("method", method).
				Str("path", path).
				Int("status", status).
				Str("response", w.body.String()).
				Dur("duration", duration).
				Msg("Handled request")
		}

		// Log errors if any
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				log.Error().
					Err(e.Err).
					Msg("Request processing error")
			}
		}
	}
}
