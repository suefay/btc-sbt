package basics

import (
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

const (
	// default tx version
	TxVersion = 2
)

const (
	// witness size for P2TR in bytes
	P2TRWitnessSize = 64

	// witness size for P2WPKH in bytes
	P2WPKHWitnessSize = 72 + 33

	// signature script size for P2SH-P2WPKH in bytes
	NestedSegWitSigScriptSize = 1 + 1 + 1 + 20

	// signature script size for P2PKH in bytes
	P2PKHSigScriptSize = 1 + 72 + 1 + 33
)

var (
	// null output script without data
	NULL_OUTPUT_SCRIPT = []byte{txscript.OP_RETURN}

	// null output
	NULL_OUTPUT = wire.NewTxOut(0, NULL_OUTPUT_SCRIPT)

	// minimum padding script without OP_RETURN
	MIN_PADDING_SCRIPT = [MIN_PADDING_SCRIPT_BYTES - 1]byte{txscript.OP_0}

	// null output script with padding
	PADDED_NULL_OUTPUT_SCRIPT = append(NULL_OUTPUT_SCRIPT, MIN_PADDING_SCRIPT[:]...)
)
