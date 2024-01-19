package params

import (
	"github.com/gin-gonic/gin"

	"github.com/btcsuite/btcd/chaincfg"
)

// Validator defines the validator interface for params validation
type Validator interface {
	Validate(c *gin.Context, netParams *chaincfg.Params) error
}
