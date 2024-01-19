package initiator

import (
	"github.com/btcsuite/btcd/txscript"

	"btc-sbt/protocol"
	"btc-sbt/stacks/basics"
)

const (
	// Default output value for issue operation when OP_RETURN not used
	DEFAULT_ISSUE_OUTPUT_VALUE = 546
)

var (
	// Default output script for mint operation
	DEFAULT_MINT_OUTPUT_SCRIPT = basics.PADDED_NULL_OUTPUT_SCRIPT

	// Output script for mint operation with the protocol name attached
	MINT_OUTPUT_SCRIPT_WITH_PROTOCOL, _ = txscript.NullDataScript([]byte(protocol.PROTOCOL_NAME))
)
