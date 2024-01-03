package basics

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// GetTxVirtualSize gets the virtual size of the given tx.
// Assume that tx.TxIn corresponds to the given utxos if the tx is unsigned
func GetTxVirtualSize(tx *wire.MsgTx, utxos []*UTXO, signed bool) int64 {
	if signed {
		return mempool.GetTxVirtualSize(btcutil.NewTx(tx))
	}

	newTx := tx.Copy()

	for i, txIn := range newTx.TxIn {
		var dummySigScript []byte
		var dummyWitness []byte

		switch txscript.GetScriptClass(utxos[i].PkScript) {
		case txscript.WitnessV1TaprootTy:
			dummyWitness = make([]byte, P2TRWitnessSize)

		case txscript.WitnessV0PubKeyHashTy:
			dummyWitness = make([]byte, P2WPKHWitnessSize)

		case txscript.ScriptHashTy:
			dummySigScript = make([]byte, NestedSegWitSigScriptSize)
			dummyWitness = make([]byte, P2WPKHWitnessSize)

		case txscript.PubKeyHashTy:
			dummySigScript = make([]byte, P2PKHSigScriptSize)

		default:
		}

		txIn.SignatureScript = dummySigScript
		txIn.Witness = wire.TxWitness{dummyWitness}
	}

	return mempool.GetTxVirtualSize(btcutil.NewTx(newTx))
}

// SerializeTx serializes the given tx
func SerializeTx(tx *wire.MsgTx) ([]byte, error) {
	var buf bytes.Buffer
	if err := tx.Serialize(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GetTotalOutputValue gets the total value of the given tx outs
func GetTotalOutputValue(txOuts []*wire.TxOut) int64 {
	var totalValue int64

	for _, out := range txOuts {
		totalValue += out.Value
	}

	return totalValue
}

// BuildTransaction builds an unsigned tx from the given params.
func BuildTransaction(utxos []*UTXO, txOuts []*wire.TxOut, paymentUtxos []*UTXO, changeAddress btcutil.Address, feeRate int64, netParams *chaincfg.Params) (*wire.MsgTx, []*UTXO, error) {
	tx := wire.NewMsgTx(TxVersion)

	inAmount := int64(0)
	outAmount := int64(0)

	for _, utxo := range utxos {
		AddUtxoToTx(tx, utxo)
		inAmount += utxo.Value
	}

	for _, txOut := range txOuts {
		tx.AddTxOut(txOut)
		outAmount += txOut.Value
	}

	changePkScript, err := txscript.PayToAddrScript(changeAddress)
	if err != nil {
		return nil, nil, err
	}

	changeOut := wire.NewTxOut(0, changePkScript)

	selectedPaymentUtxos, err := AddPaymentUtxosToTx(tx, utxos, inAmount-outAmount, paymentUtxos, changeOut, feeRate, netParams)
	if err != nil {
		return nil, nil, err
	}

	if err := ValidateTxStandardBasic(tx, netParams); err != nil {
		return nil, nil, err
	}

	return tx, selectedPaymentUtxos, nil
}

// AddPaymentUtxosToTx adds the payment utxos to the tx
func AddPaymentUtxosToTx(tx *wire.MsgTx, utxos []*UTXO, inOutdiff int64, paymentUtxos []*UTXO, changeOut *wire.TxOut, feeRate int64, netParams *chaincfg.Params) ([]*UTXO, error) {
	selectedPaymentUtxos := make([]*UTXO, 0)
	paymentValue := int64(0)

	sort.Slice(paymentUtxos, func(i, j int) bool {
		return paymentUtxos[i].Value > paymentUtxos[j].Value
	})

	for _, utxo := range paymentUtxos {
		AddUtxoToTx(tx, utxo)
		tx.AddTxOut(changeOut)

		utxos = append(utxos, utxo)
		selectedPaymentUtxos = append(selectedPaymentUtxos, utxo)

		paymentValue += utxo.Value
		fee := GetTxVirtualSize(tx, utxos, false) * feeRate

		changeValue := paymentValue + inOutdiff - fee
		if changeValue > 0 {
			tx.TxOut[len(tx.TxOut)-1].Value = changeValue
			if IsDust(tx.TxOut[len(tx.TxOut)-1], netParams) {
				tx.TxOut = tx.TxOut[0 : len(tx.TxOut)-1]
			}

			return selectedPaymentUtxos, nil
		} else {
			tx.TxOut = tx.TxOut[0 : len(tx.TxOut)-1]

			if changeValue == 0 {
				return selectedPaymentUtxos, nil
			}

			if changeValue < 0 {
				feeWithoutChange := GetTxVirtualSize(tx, utxos, false) * feeRate
				if paymentValue+inOutdiff-feeWithoutChange >= 0 {
					return selectedPaymentUtxos, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("insufficient utxos")
}
