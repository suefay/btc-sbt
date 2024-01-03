package protocol

import (
	"encoding/json"
	"strings"

	"btc-sbt/utils"
)

// OpType represents the protocol operation type
type OpType uint8

const (
	OP_UNKNOWN OpType = iota
	OP_ISSUE          // issue operation
	OP_MINT           // mint operation
)

// FromStringToOp converts the given string to OpType
func FromStringToOpType(str string) OpType {
	switch strings.ToLower(str) {
	case OP_ISSUE_TYPE_NAME:
		return OP_ISSUE

	case OP_MINT_TYPE_NAME:
		return OP_MINT

	default:
		return OP_UNKNOWN
	}
}

// String implements fmt.Stringer
func (ot OpType) String() string {
	switch ot {
	case OP_ISSUE:
		return OP_ISSUE_TYPE_NAME

	case OP_MINT:
		return OP_MINT_TYPE_NAME

	default:
		return ""
	}
}

// MarshalJSON implements json.Marshaler
func (ot OpType) MarshalJSON() ([]byte, error) {
	return json.Marshal(ot.String())
}

// UnmarshalJSON implements json.Unmarshaler
func (ot *OpType) UnmarshalJSON(data []byte) error {
	*ot = FromStringToOpType(string(data))

	return nil
}

// Operation abstracts the protocol operation
type Operation interface {
	Type() OpType
	Validate() error
	Marshal() ([]byte, error)
	Unmarshal(data []byte) error
}

// IssueOperation defines the payload struct for the issue operatioin
type IssueOperation struct {
	Op              OpType `json:"op"`               // operation type
	Symbol          string `json:"symbol"`           // unique symbol
	MaxSupply       uint64 `json:"max"`              // maximum supply
	AuthorityPubKey string `json:"authpk,omitempty"` // public key of the issuing authority
	EndBlockHeight  int64  `json:"end,omitempty"`    // end block height for mint
	Metadata        string `json:"meta,omitempty"`   // top level metadata
}

// NewIssueOperation creates an IssueOperation instance
func NewIssueOperation(symbol string, maxSupply uint64, authPK string, endBlockHeight int64, metadata string) *IssueOperation {
	return &IssueOperation{
		Op:              OP_ISSUE,
		Symbol:          symbol,
		MaxSupply:       maxSupply,
		AuthorityPubKey: authPK,
		EndBlockHeight:  endBlockHeight,
		Metadata:        metadata,
	}
}

// Type implements Operation.Type
func (op IssueOperation) Type() OpType {
	return OP_ISSUE
}

// Validate validates the issue operation
func (op *IssueOperation) Validate() error {
	if err := ValidateSymbol(op.Symbol); err != nil {
		return err
	}

	if len(op.AuthorityPubKey) > 0 {
		if err := ValidatePubKey(op.AuthorityPubKey); err != nil {
			return err
		}
	}

	if len(op.Metadata) > 0 {
		if err := ValidateMetadata(op.Metadata); err != nil {
			return err
		}
	}

	return nil
}

// Marshal marshals the IssueOperation
func (op *IssueOperation) Marshal() ([]byte, error) {
	return json.Marshal(op)
}

// Unmarshal unmarshals the given data to the IssueOperation struct
func (op *IssueOperation) Unmarshal(data []byte) error {
	return json.Unmarshal(data, op)
}

// MintOperation defines the payload struct for the mint operation
type MintOperation struct {
	Op                 OpType `json:"op"`                // operation type
	Symbol             string `json:"symbol"`            // unique symbol
	Owner              string `json:"owner"`             // token owner
	AuthoritySignature string `json:"authsig,omitempty"` // signature of the issuing authority
	Metadata           string `json:"meta,omitempty"`    // metadata per token
}

// NewMintOperation creates a MintOperation instance
func NewMintOperation(symbol string, owner string, authSig string, metadata string) *MintOperation {
	return &MintOperation{
		Op:                 OP_MINT,
		Symbol:             symbol,
		Owner:              owner,
		AuthoritySignature: authSig,
		Metadata:           metadata,
	}
}

// Type implements Operation.Type
func (op MintOperation) Type() OpType {
	return OP_MINT
}

// Validate validates the issue operation
func (op *MintOperation) Validate() error {
	if err := ValidateSymbol(op.Symbol); err != nil {
		return err
	}

	if err := ValidateAddress(op.Owner); err != nil {
		return err
	}

	if len(op.AuthoritySignature) > 0 {
		if err := ValidateSignature(op.AuthoritySignature); err != nil {
			return err
		}
	}

	if len(op.Metadata) > 0 {
		if err := ValidateMetadata(op.Metadata); err != nil {
			return err
		}
	}

	return nil
}

// Marshal marshals the MintOperation
func (op *MintOperation) Marshal() ([]byte, error) {
	return json.Marshal(op)
}

// Unmarshal unmarshals the given data to the MintOperation struct
func (op *MintOperation) Unmarshal(data []byte) error {
	return json.Unmarshal(data, op)
}

// Hash returns the sha256 hash of the MintOperation excluding AuthoritySignature
func (op MintOperation) Hash() ([]byte, error) {
	op.AuthoritySignature = ""

	bz, err := op.Marshal()
	if err != nil {
		return nil, err
	}

	return utils.SHA256(bz), nil
}
