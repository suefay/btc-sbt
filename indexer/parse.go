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
	for _, in := range tx.TxIn {
		if basics.IsTapscriptWitness(in.Witness) {
			envelope := i.Parser.GetEnvelope(in.Witness[1])
			if envelope != nil {
				ops := i.Parser.GetOps(envelope.Payload)
				if len(ops) > 0 {
					err := i.onBTCSBTProtocol(ops, tx, blockHeight, blockHash, txIndex)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// onBTCSBTProtocol performs the corresponding handling for the given protocol operations
func (i *Indexer) onBTCSBTProtocol(ops []protocol.Operation, tx *wire.MsgTx, blockHeight int64, blockHash chainhash.Hash, txIndex int) error {
	i.Logger.Infof("protocol ops found, block: %d, tx: %s", blockHeight, tx.TxHash())

	context := i.buildSMContext(blockHeight, blockHash, txIndex, tx)

	return i.StateMachine.HandleOps(context, ops)
}

// buildSMContext builds the execution context for the state machine
func (i *Indexer) buildSMContext(blockHeight int64, blockHash chainhash.Hash, txIndex int, tx *wire.MsgTx) *sm.Context {
	var firstTxOutAddr = ""

	addr, err := basics.GetAddressFromPkScript(tx.TxOut[0].PkScript, i.NetParams)
	if err == nil {
		firstTxOutAddr = addr.EncodeAddress()
	}

	return sm.NewContext(blockHeight, blockHash, txIndex, tx, firstTxOutAddr)
}
