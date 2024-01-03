package statemachine

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

// Context represents the context for handling the protocol operations
type Context struct {
	BlockHeight       int64          // block height
	BlockHash         chainhash.Hash // block hash
	TxIndex           int            // tx index
	Tx                *wire.MsgTx    // tx
	FirstTxOutAddress string         // address of the first tx out
}

// NewContext creates a new Context instance
func NewContext(blockHeight int64, blockHash chainhash.Hash, txIndex int, tx *wire.MsgTx, firstTxOutAddress string) *Context {
	return &Context{
		BlockHeight:       blockHeight,
		BlockHash:         blockHash,
		TxIndex:           txIndex,
		Tx:                tx,
		FirstTxOutAddress: firstTxOutAddress,
	}
}
