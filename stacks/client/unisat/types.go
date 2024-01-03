package unisat

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/tidwall/gjson"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"btc-sbt/stacks/basics"
)

const UNISAT_WALLET_CLIENT = "UniSat Wallet"
const UNISAT_WALLET_VERSION = "v1.1.33"

const STATUS_SUCCESS = "1"

// GetUtxosResponse defines the response struct for the utxos request
type GetUtxosResponse struct {
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Result  gjson.Result `json:"result"`
}

// GetUtxos gets the utxos from the response
func (r GetUtxosResponse) GetUtxos() ([]*basics.UTXO, error) {
	if r.Status == STATUS_SUCCESS {
		utxos := make([]*basics.UTXO, 0)

		for _, result := range r.Result.Array() {
			hash, err := chainhash.NewHashFromStr(result.Get("txId").String())
			if err != nil {
				return nil, fmt.Errorf("failed to query utxos: invalid tx hash")
			}

			scriptPk, err := hex.DecodeString(result.Get("scriptPk").String())
			if err != nil {
				return nil, fmt.Errorf("failed to query utxos: invalid pk script")
			}

			utxo := basics.NewUTXO(hash, uint32(result.Get("outputIndex").Uint()), result.Get("satoshis").Int(), scriptPk)
			utxos = append(utxos, utxo)
		}

		return utxos, nil
	}

	return nil, fmt.Errorf("failed to query utxos, err: %v", r.Message[0:20])
}

// UnmarshalJSON unmarshals the given data to the GetUtxosResponse struct
func (r *GetUtxosResponse) UnmarshalJSON(data []byte) error {
	if !gjson.ValidBytes(data) {
		return errors.New("invalid JSON")
	}

	json := gjson.ParseBytes(data)

	r.Status = json.Get("status").String()
	r.Message = json.Get("message").String()
	r.Result = json.Get("result")

	return nil
}
