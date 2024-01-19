package initiator

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"

	"btc-sbt/protocol"
	"btc-sbt/stacks/taproot/inscriber"
)

// Initiate executes the given protocol operation
func (i *Initiator) Initiate(key *secp256k1.PrivateKey, addr btcutil.Address, op protocol.Operation) (*chainhash.Hash, *chainhash.Hash, error) {
	if err := op.Validate(i.NetParams); err != nil {
		return nil, nil, err
	}

	envelope, err := GetEnvelopeFromOp(op)
	if err != nil {
		return nil, nil, err
	}

	utxos, err := i.UnisatClient.GetBTCUtxos(addr.EncodeAddress())
	if err != nil {
		return nil, nil, err
	}

	if len(utxos) == 0 {
		return nil, nil, fmt.Errorf("no utxos available for address: %s", addr)
	}

	txOut, err := GetTxOutFromOp(op, addr)
	if err != nil {
		return nil, nil, err
	}

	inscriber := inscriber.NewInscriber(i.RPCClient, i.NetParams)

	return inscriber.Inscribe(key, addr, utxos, envelope, []*wire.TxOut{txOut}, i.Config.FeeRate)
}
