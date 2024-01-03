package initiator

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"

	"btc-sbt/protocol"
	"btc-sbt/stacks/basics"
	"btc-sbt/stacks/taproot/inscriber"
)

// Initiate executes the given protocol operation
func (i *Initiator) Initiate(key *secp256k1.PrivateKey, op protocol.Operation) (*chainhash.Hash, *chainhash.Hash, error) {
	if err := op.Validate(); err != nil {
		return nil, nil, err
	}

	envelope, err := GetEnvelopeFromOp(op)
	if err != nil {
		return nil, nil, err
	}

	senderAddress, err := basics.GetTaprootAddress(key, i.NetParams)
	if err != nil {
		return nil, nil, err
	}

	senderUtxos, err := i.UnisatClient.GetAddressUtxos(senderAddress.EncodeAddress())
	if err != nil {
		return nil, nil, err
	}

	if len(senderUtxos) == 0 {
		return nil, nil, fmt.Errorf("no utxos available for address: %s", senderAddress)
	}

	txOut := GetTxOutFromOp(op, senderUtxos[0].PkScript)

	inscriber := inscriber.NewInscriber(i.RPCClient, i.NetParams)

	return inscriber.Inscribe(key, senderAddress, senderUtxos, envelope, []*wire.TxOut{txOut}, i.Config.FeeRate)
}
