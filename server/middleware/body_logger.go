package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// BodyLogger defines a logger middleware for the request body
func BodyLogger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, _ := c.Get(ReqBodyContextKey)
		if body != nil {
			bz := body.([]byte)
			logger.Infof("request body: %s", string(bz))
		}
	}
}
