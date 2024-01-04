package basics

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

const (
	// minimum relay fee for transactions in sat/kvB
	MIN_RELAY_TX_FEE = int64(1000)

	// minimum size for the btc transaction excluding witness
	MIN_STANDARD_TX_NOWITNESS_SIZE = 65

	// minimum padding bytes number in output script to satisfy the minimum tx size requirement
	MIN_PADDING_SCRIPT_BYTES = 22 // TODO: relax to 5
)

// ValidateTxStandardBasic validates the given tx against two simple standards.
// 1. dust output
// 2. minimum tx size
func ValidateTxStandardBasic(tx *wire.MsgTx, netParams *chaincfg.Params) error {
	if err := CheckDust(tx.TxOut, netParams); err != nil {
		return err
	}

	return CheckTxSize(tx)
}

// CheckDust checks if any given txout is dust output
func CheckDust(txOuts []*wire.TxOut, netParams *chaincfg.Params) error {
	for i, out := range txOuts {
		if IsDust(out, netParams) {
			return fmt.Errorf("dust output %d, value: %d", i, out.Value)
		}
	}

	return nil
}

// CheckTxSize checks if the size of the given tx conforms to the standard
func CheckTxSize(tx *wire.MsgTx) error {
	size := tx.SerializeSizeStripped()
	if size < MIN_STANDARD_TX_NOWITNESS_SIZE {
		return fmt.Errorf("tx size is small: %d, at least %d required", size, MIN_STANDARD_TX_NOWITNESS_SIZE)
	}

	return nil
}

// IsDust returns true if the given output is dust, false otherwise
func IsDust(txOut *wire.TxOut, netParams *chaincfg.Params) bool {
	if netParams.RelayNonStdTxs || txscript.IsUnspendable(txOut.PkScript) {
		return false
	}

	return mempool.IsDust(txOut, btcutil.Amount(MIN_RELAY_TX_FEE))
}
