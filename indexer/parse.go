package indexer

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"

	"btc-sbt/protocol"
	"btc-sbt/stacks/basics"
	sm "btc-sbt/statemachine"
)

// parseBTCSBTProtocol parses the potential BTC-SBT protocol data in the given block
func (i *Indexer) parseBTCSBTProtocol(blockHeight int64, block *wire.MsgBlock) error {
	for idx, tx := range block.Transactions {
		if err := i.parseBTCSBTProtocolPerTx(tx, blockHeight, block.BlockHash(), idx); err != nil {
			return err
		}
	}

	return nil
}

// parseBTCSBTProtocolPerTx parses the potential BTC-SBT protocol data in the given tx
func (i *Indexer) parseBTCSBTProtocolPerTx(tx *wire.MsgTx, blockHeight int64, blockHash chainhash.Hash, txIndex int) error {
	parsedOps := make([]protocol.Operation, 0)

	for _, in := range tx.TxIn {
		ops := i.parseBTCSBTProtocolFromWitness(in.Witness)
		if len(ops) > 0 {
			parsedOps = append(parsedOps, ops...)
			if len(parsedOps) >= protocol.BULK_OPERATION_COUNT_PER_TX {
				parsedOps = parsedOps[0:protocol.BULK_OPERATION_COUNT_PER_TX]

				break
			}
		}
	}

	if len(parsedOps) > 0 {
		return i.onBTCSBTProtocol(parsedOps, tx, blockHeight, blockHash, txIndex)
	}

	return nil
}

// parseBTCSBTProtocolFromWitness parses the potential BTC-SBT protocol data from the given witness
func (i *Indexer) parseBTCSBTProtocolFromWitness(witness wire.TxWitness) []protocol.Operation {
	if basics.IsTapscriptWitness(witness) {
		envelope := i.Parser.GetEnvelope(witness[1])
		if envelope != nil {
			return i.Parser.GetOps(envelope.Payload)
		}
	}

	return nil
}

// onBTCSBTProtocol performs the corresponding handling for the given protocol operations
func (i *Indexer) onBTCSBTProtocol(ops protocol.Operations, tx *wire.MsgTx, blockHeight int64, blockHash chainhash.Hash, txIndex int) error {
	i.Logger.Infof("protocol ops found, block: %d, tx: %s", blockHeight, tx.TxHash())

	context := i.buildSMContext(blockHeight, blockHash, txIndex, tx, ops.ContainIssue())

	return i.StateMachine.HandleOps(context, ops)
}

// buildSMContext builds the execution context for the state machine
func (i *Indexer) buildSMContext(blockHeight int64, blockHash chainhash.Hash, txIndex int, tx *wire.MsgTx, containIssueOp bool) *sm.Context {
	opOutAddr := ""

	if containIssueOp {
		opOutAddr = i.Parser.ParseIssuerAddress(tx)
	}

	return sm.NewContext(blockHeight, blockHash, txIndex, tx, opOutAddr)
}
