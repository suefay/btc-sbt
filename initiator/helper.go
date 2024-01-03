package initiator

import (
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
func GetTxOutFromOp(op protocol.Operation, pkScript []byte) *wire.TxOut {
	switch op.Type() {
	case protocol.OP_ISSUE:
		return wire.NewTxOut(DEFAULT_ISSUE_OUTPUT_VALUE, pkScript)

	case protocol.OP_MINT:
		return wire.NewTxOut(0, DEFAULT_MINT_OUTPUT_SCRIPT)

	default:
		return nil
	}
}
