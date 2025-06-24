package pkg

import (
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		statusCode := c.Writer.Status()
		path := c.FullPath()
		method := c.Request.Method

		log.WithFields(log.Fields{
			"status":   statusCode,
			"method":   method,
			"path":     path,
			"duration": duration.String(),
		}).Info("HTTP request processed")
	}
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.DebugLevel)
}
