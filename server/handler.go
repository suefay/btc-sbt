package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"btc-sbt/server/params"
)

// GetAllSBTs queries all the SBTs
func (srv *APIService) GetAllSBTs(c *gin.Context) {
	collections, err := srv.APIBackend.GetAllSBTs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": fmt.Sprintf("%v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "result": collections})
}

// GetSBTs queries the SBTs by the given symbol
func (srv *APIService) GetSBTs(c *gin.Context) {
	var p params.GetSBTsParams
	if err := c.ShouldBindUri(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": fmt.Sprintf("invalid params: %v", err)})
		return
	}

	if err := p.Validate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": fmt.Sprintf("invalid params: %v", err)})
		return
	}

	sbts, err := srv.APIBackend.GetSBTs(p.Symbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": fmt.Sprintf("%v", err)})
		return
	}

	if sbts == nil {
		c.JSON(http.StatusNotFound, gin.H{"status": false, "error": fmt.Sprintf("SBTs does not exist: %s", p.Symbol)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "result": sbts})
}

// GetSBT queries the SBT token by the given symbol and token id
func (srv *APIService) GetSBT(c *gin.Context) {
	var p params.GetSBTParams
	if err := c.ShouldBindQuery(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": fmt.Sprintf("invalid params: %v", err)})
		return
	}

	if err := p.Validate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": fmt.Sprintf("invalid params: %v", err)})
		return
	}

	sbt, err := srv.APIBackend.GetSBT(p.Symbol, p.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": fmt.Sprintf("%v", err)})
		return
	}

	if sbt == nil {
		c.JSON(http.StatusNotFound, gin.H{"status": false, "error": fmt.Sprintf("SBT does not exist, symbol: %s, id: %d", p.Symbol, p.Id)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "result": sbt})
}

// GetOwnedSBTsWrapper dispatches execution to the GetOwnedSBTs handler or GetOwnedSBT handler according to request params
func (srv *APIService) GetOwnedSBTsWrapper(c *gin.Context) {
	var p params.GetOwnedSBTsWrapperParams
	if err := c.ShouldBindQuery(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": fmt.Sprintf("invalid params: %v", err)})
		return
	}

	if !p.SymbolExists() {
		srv.GetOwnedSBTs(c)
	} else {
		srv.GetOwnedSBT(c)
	}
}

// GetOwnedSBTs queries the SBT tokens owned by the given address
func (srv *APIService) GetOwnedSBTs(c *gin.Context) {
	var p params.GetOwnedSBTsParams
	if err := c.ShouldBindUri(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": fmt.Sprintf("invalid params: %v", err)})
		return
	}

	if err := p.Validate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": fmt.Sprintf("invalid params: %v", err)})
		return
	}

	sbts, err := srv.APIBackend.GetOwnedSBTs(p.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": fmt.Sprintf("%v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "result": sbts})
}

// GetOwnedSBT queries the specified SBT token owned by the given address
func (srv *APIService) GetOwnedSBT(c *gin.Context) {
	var p params.GetOwnedSBTParams

	if err := c.ShouldBindUri(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": fmt.Sprintf("invalid params: %v", err)})
		return
	}

	if err := c.ShouldBindQuery(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": fmt.Sprintf("invalid params: %v", err)})
		return
	}

	if err := p.Validate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": fmt.Sprintf("invalid params: %v", err)})
		return
	}

	sbt, err := srv.APIBackend.GetOwnedSBT(p.Address, p.Symbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": fmt.Sprintf("%v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "result": sbt})
}

// Status returns the current status of the indexer
func (srv *APIService) Status(c *gin.Context) {
	res, err := srv.APIBackend.GetStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": fmt.Sprintf("%v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "result": res})
}
