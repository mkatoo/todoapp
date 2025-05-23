package accesslog

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func AccessLogMiddleware() gin.HandlerFunc {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return func(c *gin.Context) {
		start := time.Now()

		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ := io.ReadAll(tee)
		c.Request.Body = io.NopCloser(&buf)

		// Process request
		c.Next()

		end := time.Now()

		logger.Info("Request",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Duration("duration", time.Duration(end.Sub(start))),
			slog.Int("status", c.Writer.Status()),
			slog.String("body", string(body)),
		)
	}
}
