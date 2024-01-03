package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// WithBody attaches the request body to the context for underlying handling by other middlewares
func WithBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		body := c.Request.Body
		if body != nil && body != http.NoBody {
			bz, _ := io.ReadAll(body)
			if bz != nil {
				c.Request.Body = io.NopCloser(bytes.NewReader(bz))
				c.Set(ReqBodyContextKey, bz)
			}
		}
	}
}
