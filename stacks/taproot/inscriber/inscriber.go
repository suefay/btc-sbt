package inscriber

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"

	"btc-sbt/stacks/basics"
	"btc-sbt/stacks/taproot"
)

// Inscriber defines the inscriber struct
type Inscriber struct {
	rpcClient *rpcclient.Client

	netParams *chaincfg.Params
}

// NewInscriber creates an Inscriber instance
func NewInscriber(rpcClient *rpcclient.Client, netParams *chaincfg.Params) *Inscriber {
	return &Inscriber{
		rpcClient,
		netParams,
	}
}

// Inscribe performs the inscribing process which consists of two phases named commit and reveal
func (i *Inscriber) Inscribe(commitKey *secp256k1.PrivateKey, commitAddress *btcutil.AddressTaproot, commitUtxos []*basics.UTXO, envelope []byte, revealTxOuts []*wire.TxOut, feeRate int64) (*chainhash.Hash, *chainhash.Hash, error) {
	commitOutWIF, commitOutAddress, err := taproot.GenerateTapscriptCommitOutAddress(envelope, i.netParams)
	if err != nil {
		return nil, nil, err
	}

	commitOutPkScript, err := txscript.PayToAddrScript(commitOutAddress)
	if err != nil {
		return nil, nil, err
	}

	script, err := taproot.BuildTapscript(commitOutWIF.PrivKey.PubKey(), envelope)
	if err != nil {
		return nil, nil, err
	}

	revealTx, err := i.buildDummyRevealTx(commitOutWIF.PrivKey.PubKey(), script, revealTxOuts)
	if err != nil {
		return nil, nil, err
	}

	commitOutValue := basics.GetTxVirtualSize(revealTx, nil, true)*feeRate + basics.GetTotalOutputValue(revealTxOuts)
	commitTx, err := i.buildCommitTx(commitKey, commitAddress, commitUtxos, commitOutValue, commitOutPkScript, feeRate)
	if err != nil {
		return nil, nil, err
	}

	commitTxHash := commitTx.TxHash()
	commitTxOutPoint := wire.NewOutPoint(&commitTxHash, 0)
	commitTxOut := wire.NewTxOut(commitOutValue, commitOutPkScript)

	if err := i.populateDummyRevealTx(commitOutWIF.PrivKey, revealTx, *commitTxOutPoint, commitTxOut); err != nil {
		return nil, nil, err
	}

	_, err = i.rpcClient.SendRawTransaction(commitTx, false)
	if err != nil {
		return nil, nil, err
	}

	revealTxHash, err := i.rpcClient.SendRawTransaction(revealTx, false)
	if err != nil {
		return nil, nil, err
	}

	return &commitTxHash, revealTxHash, nil
}

func (i *Inscriber) buildCommitTx(key *secp256k1.PrivateKey, commitAddress btcutil.Address, utxos []*basics.UTXO, outValue int64, outPkScript []byte, feeRate int64) (*wire.MsgTx, error) {
	commitOut := wire.NewTxOut(outValue, outPkScript)

	tx, utxos, err := basics.BuildTransaction(nil, []*wire.TxOut{commitOut}, utxos, commitAddress, feeRate, i.netParams)
	if err != nil {
		return nil, fmt.Errorf("failed to build commit tx: %v", err)
	}

	if err := i.signCommitTx(key, tx, utxos, txscript.SigHashDefault); err != nil {
		return nil, fmt.Errorf("failed to build commit tx: %v", err)
	}

	return tx, nil
}

func (i *Inscriber) buildDummyRevealTx(pubKey *secp256k1.PublicKey, script []byte, txOuts []*wire.TxOut) (*wire.MsgTx, error) {
	tx := wire.NewMsgTx(basics.TxVersion)

	tx.AddTxIn(&wire.TxIn{})

	if len(txOuts) == 0 {
		txOuts = append(txOuts, wire.NewTxOut(0, basics.PADDED_NULL_OUTPUT_SCRIPT))
	}

	for _, out := range txOuts {
		tx.AddTxOut(out)
	}

	if err := basics.ValidateTxStandardBasic(tx, i.netParams); err != nil {
		return nil, fmt.Errorf("failed to build reveal tx: %v", err)
	}

	witness, err := taproot.GetTapscriptWitnessNoSignature(pubKey, script, txscript.SigHashDefault)
	if err != nil {
		return nil, fmt.Errorf("failed to build reveal tx: %v", err)
	}

	tx.TxIn[0].Witness = witness

	return tx, nil
}

func (i *Inscriber) populateDummyRevealTx(key *secp256k1.PrivateKey, revealTx *wire.MsgTx, commitTxOutPoint wire.OutPoint, commitTxOut *wire.TxOut) error {
	revealTx.TxIn[0].PreviousOutPoint = commitTxOutPoint

	if err := i.signRevealTx(key, revealTx, commitTxOut, txscript.SigHashDefault); err != nil {
		return fmt.Errorf("failed to populate reveal tx: %v", err)
	}

	return nil
}

func (i *Inscriber) signCommitTx(key *secp256k1.PrivateKey, tx *wire.MsgTx, utxos []*basics.UTXO, hashType txscript.SigHashType) error {
	return taproot.SignTaprootTransaction(key, tx, utxos, hashType)
}

func (i *Inscriber) signRevealTx(key *secp256k1.PrivateKey, tx *wire.MsgTx, commitTxOut *wire.TxOut, hashType txscript.SigHashType) error {
	revealScript := tx.TxIn[0].Witness[1]
	utxo := &basics.UTXO{Value: commitTxOut.Value, PkScript: commitTxOut.PkScript}

	signature, err := taproot.SignTapscript(key, tx, []*basics.UTXO{utxo}, 0, revealScript, hashType)
	if err != nil {
		return err
	}

	tx.TxIn[0].Witness[0] = signature

	return nil
}
