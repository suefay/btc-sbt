package basics

import (
	"encoding/hex"
	"errors"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// GetAddressFromPkScript gets the btc address from the given pk script and net params
func GetAddressFromPkScript(pkScript []byte, netParams *chaincfg.Params) (btcutil.Address, error) {
	ps, err := txscript.ParsePkScript(pkScript)
	if err != nil {
		return nil, err
	}

	return ps.Address(netParams)
}

// IsTapscriptWitness returns true if the given witness is tapscript witness, false otherwise
func IsTapscriptWitness(witness wire.TxWitness) bool {
	if len(witness) < 3 {
		return false
	}

	_, err := txscript.ParseControlBlock(witness[2])
	return err == nil
}

// GetRedeemScriptForNestedSegWit gets the redeem script for the P2SH-P2WPKH address
func GetRedeemScriptForNestedSegWit(pubKey string, netParams *chaincfg.Params) ([]byte, error) {
	if len(pubKey) == 0 {
		return nil, errors.New("empty public key")
	}

	pubKeyBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		return nil, err
	}

	p2wpkh, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pubKeyBytes), netParams)
	if err != nil {
		return nil, err
	}

	return txscript.PayToAddrScript(p2wpkh)
}

// AddLargeDataToScript adds a large block of data to the script
func AddLargeDataToScript(scriptBuilder *txscript.ScriptBuilder, data []byte) {
	maxChunkSize := txscript.MaxScriptElementSize

	dataSize := len(data)
	for i := 0; i < dataSize; i += maxChunkSize {
		end := i + maxChunkSize
		if end > dataSize {
			end = dataSize
		}

		scriptBuilder.AddFullData(data[i:end])
	}
}
