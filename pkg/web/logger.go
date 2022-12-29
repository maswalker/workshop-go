package web

import (
	"time"

	"go-seven/pkg/log"

	"github.com/gin-gonic/gin"
)

func Logger(logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)

		fields := []log.Field{
			log.Int("status", c.Writer.Status()),
			log.String("method", c.Request.Method),
			log.String("path", path),
			log.String("query", query),
			log.String("ip", c.ClientIP()),
			log.String("user-agent", c.Request.UserAgent()),
			log.Duration("cost", cost),
		}
		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				logger.Error(e, fields...)
			}
		} else {
			logger.Info(path, fields...)
		}
	}
}
