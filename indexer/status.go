package indexer

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

// GetLastBlockHeight gets the last block height from the indexer
func (i *Indexer) GetLastBlockHeight() (int64, error) {
	return i.StateMachine.GetLastBlockHeight()
}

// GetLastBlockHeight gets the last block hash from the indexer
func (i *Indexer) GetLastBlockHash() (*chainhash.Hash, error) {
	return i.StateMachine.GetLastBlockHash()
}

// GetLatestBlockHeight gets the latest block height from the chain
func (i *Indexer) GetLatestBlockHeight() (int64, error) {
	return i.Client.GetLatestBlockHeight()
}
