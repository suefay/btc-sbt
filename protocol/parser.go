package protocol

import (
	"fmt"

	"github.com/tidwall/gjson"
)

// Parser defines the parser for the BTC-SBT protocol
type Parser struct {
}

// NewParser creates a new Parser instance
func NewParser() *Parser {
	return &Parser{}
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
