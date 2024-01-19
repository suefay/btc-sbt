package indexer

import (
	"github.com/btcsuite/btcd/chaincfg"

	"btc-sbt/types"
)

// GetSBTs queries all the SBTs
func (i *Indexer) GetAllSBTs() ([]*types.SBTs, error) {
	return i.StateMachine.GetAllSBTs()
}

// GetSBTs queries the SBTs by the given symbol
func (i *Indexer) GetSBTs(symbol string) (*types.SBTs, error) {
	return i.StateMachine.GetSBTs(symbol)
}

// GetSBT queries the SBT token by the given symbol and token id
func (i *Indexer) GetSBT(symbol string, id uint64) (*types.SBT, error) {
	return i.StateMachine.GetSBT(symbol, id)
}

// GetOwnedSBTs queries the SBT tokens owned by the given owner
func (i *Indexer) GetOwnedSBTs(owner string) ([]*types.CompactSBT, error) {
	return i.StateMachine.GetOwnedSBTs(owner)
}

// GetOwnedSBT queries the specified SBT token owned by the given owner
func (i *Indexer) GetOwnedSBT(owner string, symbol string) (*types.CompactSBT, error) {
	return i.StateMachine.GetOwnedSBT(owner, symbol)
}

// GetStatus returns the current status of the indexer
func (i *Indexer) GetStatus() (any, error) {
	return i.GetLastBlockHeight()
}

// GetNetParams returns the net params used for the indexer
func (i *Indexer) GetNetParams() *chaincfg.Params {
	return i.NetParams
}
