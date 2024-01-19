package initiator

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"

	"btc-sbt/protocol"
)

// GetEnvelopeFromOp gets the envelope from the given operation
func GetEnvelopeFromOp(op protocol.Operation) ([]byte, error) {
	payload, err := op.Marshal()
	if err != nil {
		return nil, err
	}

	return protocol.NewEnvelope(payload).Script()
}

// GetTxOutFromOp gets the txout from the given operation
func GetTxOutFromOp(op protocol.Operation, addr btcutil.Address) (*wire.TxOut, error) {
	switch op.Type() {
	case protocol.OP_ISSUE:
		script, err := txscript.NullDataScript([]byte(addr.EncodeAddress()))
		if err != nil {
			return nil, fmt.Errorf("failed to get tx out from operation, err: %v", err)
		}

		return wire.NewTxOut(0, script), nil

	case protocol.OP_MINT:
		return wire.NewTxOut(0, MINT_OUTPUT_SCRIPT_WITH_PROTOCOL), nil

	default:
		return nil, fmt.Errorf("unsupported operation type: %d", op.Type())
	}
}
