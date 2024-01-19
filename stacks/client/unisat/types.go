package unisat

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"btc-sbt/stacks/basics"
)

const UNISAT_WALLET_CLIENT = "UniSat Wallet"

// UTXO represents the UTXO struct for unisat api
type UTXO struct {
	TxId        string `json:"txid"`
	Vout        uint32 `json:"vout"`
	Value       int64  `json:"satoshis"`
	PkScript    string `json:"scriptPk"`
	AddressType uint8  `json:"addressType"`
	Height      int64  `json:"height"`
}

// GetBTCUtxosResponse defines the response struct for the GetBTCUtxos request
type GetBTCUtxosResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    []UTXO `json:"data"`
}

// GetUtxos gets the utxos from the response
func (r GetBTCUtxosResponse) GetUtxos() ([]*basics.UTXO, error) {
	if r.Code == 0 {
		utxos := make([]*basics.UTXO, 0)

		for _, utxo := range r.Data {
			hash, err := chainhash.NewHashFromStr(utxo.TxId)
			if err != nil {
				return nil, fmt.Errorf("invalid response: invalid tx hash")
			}

			pkScript, err := hex.DecodeString(utxo.PkScript)
			if err != nil {
				return nil, fmt.Errorf("invalid response: invalid pk script")
			}

			utxo := basics.NewUTXO(hash, utxo.Vout, utxo.Value, pkScript)
			utxos = append(utxos, utxo)
		}

		return utxos, nil
	}

	return nil, fmt.Errorf("%v", r.Message)
}

// Unmarshal unmarshals the given data to the GetBTCUtxosResponse struct
func (r *GetBTCUtxosResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}
