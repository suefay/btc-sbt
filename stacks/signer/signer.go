package signer

import (
	"btc-sbt/stacks/basics"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

func SignWitnessTransaction(key *secp256k1.PrivateKey, tx *wire.MsgTx, utxos []*basics.UTXO, hashType txscript.SigHashType) error {
	prevOutFetcher := txscript.NewMultiPrevOutFetcher(nil)

	for i, utxo := range utxos {
		prevOutFetcher.AddPrevOut(tx.TxIn[i].PreviousOutPoint, utxo.GetOutput())
	}

	for i, txIn := range tx.TxIn {
		witness, err := txscript.WitnessSignature(tx, txscript.NewTxSigHashes(tx, prevOutFetcher), i, utxos[i].Value, utxos[i].PkScript, hashType, key, true)
		if err != nil {
			return err
		}

		txIn.Witness = witness
	}

	return nil
}
