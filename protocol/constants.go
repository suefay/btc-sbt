package protocol

import (
	"github.com/btcsuite/btcd/txscript"
)

const (
	// Protocol name
	PROTOCOL_NAME = "BTC-SBT"

	// Protocol identifier
	PROTOCOL_IDENTIFIER = "sbt"
)

// Envelope related constants
var (
	ENVELOPE_HEADER              = []byte{txscript.OP_FALSE, txscript.OP_IF}
	ENVELOPE_PROTOCOL_IDENTIFIER = []byte(PROTOCOL_IDENTIFIER)
	ENVELOPE_PAYLOAD_TAG         = byte(txscript.OP_0)
	ENVELOPE_TAIL                = []byte{txscript.OP_ENDIF}
)

const (
	// Operation type names
	OP_ISSUE_TYPE_NAME = "issue"
	OP_MINT_TYPE_NAME  = "mint"

	// The issuer is identified by the first output
	ISSUER_OUTPUT_INDEX = 0

	// Bulk operation count per tx
	BULK_OPERATION_COUNT_PER_TX = 10
)

const (
	// Minimum size of the symbol in bytes
	MIN_SYMBOL_LEN = 3

	// Miximum size of the symbol in bytes
	MAX_SYMBOL_LEN = 8
)
