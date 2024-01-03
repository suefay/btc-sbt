package params

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"btc-sbt/protocol"
)

var _ Validator = (*GetSBTsParams)(nil)
var _ Validator = (*GetSBTParams)(nil)
var _ Validator = (*GetOwnedSBTsWrapperParams)(nil)
var _ Validator = (*GetOwnedSBTsParams)(nil)
var _ Validator = (*GetOwnedSBTParams)(nil)

// GetSBTsParams represents the params for the GetSBTs handler
type GetSBTsParams struct {
	Symbol string `json:"symbol" uri:"symbol"`
}

// Validate implements the Validator interface
func (p *GetSBTsParams) Validate(c *gin.Context) error {
	return protocol.ValidateSymbol(p.Symbol)
}

// GetSBTParams represents the params for the GetSBT handler
type GetSBTParams struct {
	Symbol string `json:"symbol" form:"symbol"`
	Id     uint64 `json:"id" form:"id"`
}

// Validate implements the Validator interface
func (p *GetSBTParams) Validate(c *gin.Context) error {
	if err := protocol.ValidateSymbol(p.Symbol); err != nil {
		return err
	}

	id := c.Query("id")
	if len(id) == 0 {
		return fmt.Errorf("token id missing")
	}

	return nil
}

// GetOwnedSBTsWrapperParams represents the params for the GetOwnedSBTsWrapper handler
type GetOwnedSBTsWrapperParams struct {
	Symbol string `json:"symbol" form:"symbol"`
}

// Validate implements the Validator interface
func (p *GetOwnedSBTsWrapperParams) Validate(c *gin.Context) error {
	// defer validation to concrete handler
	return nil
}

// SymbolExists returns true if the params contain non-empty `symbol`, false otherwise
func (p *GetOwnedSBTsWrapperParams) SymbolExists() bool {
	return len(p.Symbol) > 0
}

// GetOwnedSBTsParams represents the params for the GetOwnedSBTs handler
type GetOwnedSBTsParams struct {
	Address string `json:"address" uri:"address"`
}

// Validate implements the Validator interface
func (p *GetOwnedSBTsParams) Validate(c *gin.Context) error {
	return protocol.ValidateAddress(p.Address)
}

// GetOwnedSBTParams represents the params for the GetOwnedSBT handler
type GetOwnedSBTParams struct {
	Address string `json:"address" uri:"address"`
	Symbol  string `json:"symbol" form:"symbol"`
}

// Validate implements the Validator interface
func (p *GetOwnedSBTParams) Validate(c *gin.Context) error {
	if err := protocol.ValidateAddress(p.Address); err != nil {
		return err
	}

	return protocol.ValidateSymbol(p.Symbol)
}
