package protocol

import (
	"fmt"

	"github.com/tidwall/gjson"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"

	"btc-sbt/stacks/basics"
)

// Parser defines the parser for the BTC-SBT protocol
type Parser struct {
	NetParams *chaincfg.Params // net params
}

// NewParser creates a new Parser instance
func NewParser(netParams *chaincfg.Params) *Parser {
	return &Parser{
		NetParams: netParams,
	}
}

// GetEnvelope gets the BTC-SBT protocol envelope from the given witness script
func (p *Parser) GetEnvelope(witnessScript []byte) *Envelope {
	return ExtractEnvelope(witnessScript)
}

// GetOps parses the BTC-SBT protocol operations from the given payload
func (p *Parser) GetOps(payload []byte) []Operation {
	if !gjson.ValidBytes(payload) {
		return nil
	}

	json := gjson.ParseBytes(payload)

	operations := make([]Operation, 0)

	if json.IsArray() {
		for _, e := range json.Array() {
			op, err := p.getOpFromJSON(e)
			if err == nil {
				operations = append(operations, op)
			}
		}
	} else {
		op, err := p.getOpFromJSON(json)
		if err == nil {
			operations = append(operations, op)
		}
	}

	return operations
}

// ToIssueOp parses the given data to the issue operation
func (p *Parser) ToIssueOp(data []byte) (*IssueOperation, error) {
	var op IssueOperation
	err := op.Unmarshal(data)

	return &op, err
}

// ToMintOp parses the given data to the mint operation
func (p *Parser) ToMintOp(data []byte) (*MintOperation, error) {
	var op MintOperation
	err := op.Unmarshal(data)

	return &op, err
}

// ParseIssuerAddress parses the issuer address from the given tx.
// Assume that the tx contains the issue operation(s)
func (p *Parser) ParseIssuerAddress(tx *wire.MsgTx) string {
	if len(tx.TxOut) <= ISSUER_OUTPUT_INDEX {
		return ""
	}

	pkScript := tx.TxOut[ISSUER_OUTPUT_INDEX].PkScript

	if txscript.IsNullData(pkScript) {
		return p.parseIssuerAddressFromNullData(pkScript)
	}

	addr, err := basics.GetAddressFromPkScript(pkScript, p.NetParams)
	if err == nil {
		return addr.EncodeAddress()
	}

	return ""
}

// parseIssuerAddressFromNullData parses the issuer address from the given null data script
func (p *Parser) parseIssuerAddressFromNullData(nullDataScript []byte) string {
	if len(nullDataScript) <= 2 {
		return ""
	}

	// data push
	addr, err := btcutil.DecodeAddress(string(nullDataScript[2:]), p.NetParams)
	if err != nil {
		return ""
	}

	return addr.EncodeAddress()
}

// getOpFromJSON parses the BTC-SBT protocol operation from the given JSON data
func (p *Parser) getOpFromJSON(data gjson.Result) (Operation, error) {
	op := data.Get("op")
	if op.Type != gjson.String {
		return nil, fmt.Errorf("invalid op data type, string required")
	}

	switch op.Str {
	case OP_ISSUE.String():
		return p.ToIssueOp([]byte(data.Raw))

	case OP_MINT.String():
		return p.ToMintOp([]byte(data.Raw))

	default:
		return nil, fmt.Errorf("unknown op: %s", op.Str)
	}
}
