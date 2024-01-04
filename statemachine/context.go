package statemachine

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

// Context represents the context for handling the protocol operations
type Context struct {
	BlockHeight         int64          // block height
	BlockHash           chainhash.Hash // block hash
	TxIndex             int            // tx index
	Tx                  *wire.MsgTx    // tx
	OperationOutAddress string         // operation output address, i.e. issuer address if there exists the issue operation
}

// NewContext creates a new Context instance
func NewContext(blockHeight int64, blockHash chainhash.Hash, txIndex int, tx *wire.MsgTx, opOutAddr string) *Context {
	return &Context{
		BlockHeight:         blockHeight,
		BlockHash:           blockHash,
		TxIndex:             txIndex,
		Tx:                  tx,
		OperationOutAddress: opOutAddr,
	}
}
