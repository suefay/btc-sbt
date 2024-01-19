package server

import (
	"btc-sbt/types"

	"github.com/btcsuite/btcd/chaincfg"
)

// APIBackend defines the backend interface for the api server
type APIBackend interface {
	GetAllSBTs() ([]*types.SBTs, error)

	GetSBTs(symbol string) (*types.SBTs, error)
	GetSBT(symbol string, id uint64) (*types.SBT, error)

	GetOwnedSBTs(owner string) ([]*types.CompactSBT, error)
	GetOwnedSBT(owner string, symbol string) (*types.CompactSBT, error)

	GetStatus() (any, error)

	GetNetParams() *chaincfg.Params
}
