package psbt

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"

	"btc-sbt/stacks/basics"
)

// AddInputToPsbt adds the given utxo to the psbt.
// Assume that the input index is valid
func AddInputToPsbt(p *psbt.Packet, index int, utxo *basics.UTXO, rawTx *wire.MsgTx, address btcutil.Address, pubKey string, sigHashType txscript.SigHashType, netParams *chaincfg.Params) error {
	p.Inputs[index].SighashType = sigHashType

	if basics.IsSegWitAddress(address) {
		p.Inputs[index].WitnessUtxo = utxo.GetOutput()

		if basics.IsTaprootAddress(address) {
			xOnlyPubKey, err := basics.GetXOnlyPubKey(pubKey)
			if err != nil {
				return err
			}

			p.Inputs[index].TaprootInternalKey = xOnlyPubKey
		}

		return nil
	}

	if basics.IsP2SHAddress(address) {
		redeemScript, err := basics.GetRedeemScriptForNestedSegWit(pubKey, netParams)
		if err != nil {
			return err
		}

		p.Inputs[index].RedeemScript = redeemScript
	}

	p.Inputs[index].NonWitnessUtxo = rawTx

	return nil
}
