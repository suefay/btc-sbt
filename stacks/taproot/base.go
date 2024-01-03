package taproot

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"

	"btc-sbt/stacks/basics"
)

func GenerateTapscriptCommitOutAddress(envelope []byte, netParams *chaincfg.Params) (*btcutil.WIF, *btcutil.AddressTaproot, error) {
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, nil, err
	}

	wif, err := btcutil.NewWIF(privKey, netParams, true)
	if err != nil {
		return nil, nil, err
	}

	script, err := BuildTapscript(privKey.PubKey(), envelope)
	if err != nil {
		return nil, nil, err
	}

	tapHash := txscript.NewBaseTapLeaf(script).TapHash()

	commitOutAddress, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootOutputKey(privKey.PubKey(), tapHash[:])), netParams)
	if err != nil {
		return nil, nil, err
	}

	return wif, commitOutAddress, nil
}

func BuildTapscript(pubKey *btcec.PublicKey, envelope []byte) ([]byte, error) {
	scriptBuilder := txscript.NewScriptBuilder().
		AddData(schnorr.SerializePubKey(pubKey)).
		AddOp(txscript.OP_CHECKSIG)

	script, err := scriptBuilder.Script()
	if err != nil {
		return nil, err
	}

	return append(script, envelope...), nil
}

func GetTapscriptWitnessNoSignature(pubKey *btcec.PublicKey, script []byte, hashType txscript.SigHashType) (wire.TxWitness, error) {
	proof := txscript.TapscriptProof{
		TapLeaf:  txscript.NewBaseTapLeaf(script),
		RootNode: txscript.NewBaseTapLeaf(script),
	}

	controlBlock := proof.ToControlBlock(pubKey)

	controlBlockBz, err := controlBlock.ToBytes()
	if err != nil {
		return nil, err
	}

	signatureSize := schnorr.SignatureSize
	if hashType&txscript.SigHashDefault != txscript.SigHashDefault {
		signatureSize += 1
	}

	witness := make(wire.TxWitness, 3)

	witness[0] = make([]byte, signatureSize)
	witness[1] = script
	witness[2] = controlBlockBz

	return witness, nil
}

func SignTaproot(key *secp256k1.PrivateKey, tx *wire.MsgTx, utxos []*basics.UTXO, idx int, hashType txscript.SigHashType) ([][]byte, error) {
	prevOutFetcher := txscript.NewMultiPrevOutFetcher(nil)

	for i, utxo := range utxos {
		prevOutFetcher.AddPrevOut(tx.TxIn[i].PreviousOutPoint, utxo.GetOutput())
	}

	witness, err := txscript.TaprootWitnessSignature(tx, txscript.NewTxSigHashes(tx, prevOutFetcher), idx, utxos[idx].Value, utxos[idx].PkScript, hashType, key)
	if err != nil {
		return nil, err
	}

	return witness, nil
}

func SignTaprootTransaction(key *secp256k1.PrivateKey, tx *wire.MsgTx, utxos []*basics.UTXO, hashType txscript.SigHashType) error {
	prevOutFetcher := txscript.NewMultiPrevOutFetcher(nil)

	for i, utxo := range utxos {
		prevOutFetcher.AddPrevOut(tx.TxIn[i].PreviousOutPoint, utxo.GetOutput())
	}

	for i, txIn := range tx.TxIn {
		witness, err := txscript.TaprootWitnessSignature(tx, txscript.NewTxSigHashes(tx, prevOutFetcher), i, utxos[i].Value, utxos[i].PkScript, hashType, key)
		if err != nil {
			return err
		}

		txIn.Witness = witness
	}

	return nil
}

func SignTapscript(key *secp256k1.PrivateKey, tx *wire.MsgTx, utxos []*basics.UTXO, idx int, script []byte, hashType txscript.SigHashType) ([]byte, error) {
	prevOutFetcher := txscript.NewMultiPrevOutFetcher(nil)

	for i, utxo := range utxos {
		prevOutFetcher.AddPrevOut(tx.TxIn[i].PreviousOutPoint, utxo.GetOutput())
	}

	signature, err := txscript.RawTxInTapscriptSignature(tx, txscript.NewTxSigHashes(tx, prevOutFetcher), idx, utxos[idx].Value, utxos[idx].PkScript, txscript.NewBaseTapLeaf(script), hashType, key)
	if err != nil {
		return nil, err
	}

	return signature, nil
}
