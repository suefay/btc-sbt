package cmd

import (
	"os"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"

	"btc-sbt/stacks/basics"
)

// GetPrivateKeyAndTaprootAddr gets the private key and corresponding taproot address from the given file path
func GetPrivateKeyAndTaprootAddr(path string, netParam *chaincfg.Params) (*secp256k1.PrivateKey, *btcutil.AddressTaproot, error) {
	keyWIFBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	keyWIF, err := btcutil.DecodeWIF(string(keyWIFBytes))
	if err != nil {
		return nil, nil, err
	}

	taprootAddr, err := basics.GetTaprootAddress(keyWIF.PrivKey, netParam)
	if err != nil {
		return nil, nil, err
	}

	return keyWIF.PrivKey, taprootAddr, nil
}
