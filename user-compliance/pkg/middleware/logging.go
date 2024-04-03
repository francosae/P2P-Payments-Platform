package middleware

import (
	"bytes"
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
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		w := &responseWriter{body: new(bytes.Buffer), ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()

		duration := time.Since(start)
		log.Info().
			Str("method", method).
			Str("path", path).
			Int("status", c.Writer.Status()).
			Str("response", w.body.String()).
			Dur("duration", duration).
			Msg("handled request")
	}
}
