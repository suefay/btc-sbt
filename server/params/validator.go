package params

import (
	"github.com/gin-gonic/gin"
)

// Validator defines the validator interface for params validation
type Validator interface {
	Validate(c *gin.Context) error
}
