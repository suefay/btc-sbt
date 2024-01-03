package types

import (
	"encoding/json"

	"btc-sbt/protocol"
)

// SBTs defines the token struct
type SBTs struct {
	Symbol          string `json:"symbol"`              // unique symbol
	Sequence        uint64 `json:"seq"`                 // increasing sequence number, starting from 1
	MaxSupply       uint64 `json:"max_supply"`          // maximum supply
	AuthorityPubKey string `json:"auth_pk,omitempty"`   // public key of the issuing authority
	EndBlockHeight  int64  `json:"end_block,omitempty"` // end block height for mint
	Metadata        string `json:"metadata,omitempty"`  // top level metadata

	Issuer               string `json:"issuer"`       // issuer address
	BlockHeight          int64  `json:"block_height"` // issue block height
	TransactionIndex     int    `json:"tx_index"`     // issue tx index
	IssueTransactionHash string `json:"issue_tx"`     // issue tx hash

	TotalSupply uint64 `json:"total_supply,omitempty"` // current supply, i.e. total amount of minted tokens
}

// NewSBTs creates a new SBTs instance
func NewSBTs(symbol string, seq uint64, maxSupply uint64, pubKey string, endBlockHeight int64, metadata string, issuer string, blockHeight int64, txIndex int, issueTxHash string, totalSupply uint64) *SBTs {
	return &SBTs{
		Symbol:               symbol,
		Sequence:             seq,
		MaxSupply:            maxSupply,
		AuthorityPubKey:      pubKey,
		EndBlockHeight:       endBlockHeight,
		Metadata:             metadata,
		Issuer:               issuer,
		BlockHeight:          blockHeight,
		TransactionIndex:     txIndex,
		IssueTransactionHash: issueTxHash,
		TotalSupply:          totalSupply,
	}
}

// NewSBTsFromIssueOp creates an SBTs from the given issue operation
func NewSBTsFromIssueOp(op *protocol.IssueOperation) *SBTs {
	return NewSBTs(
		op.Symbol,
		0,
		op.MaxSupply,
		op.AuthorityPubKey,
		op.EndBlockHeight,
		op.Metadata,
		"",
		0,
		0,
		"",
		0,
	)
}

// RequireSignatureOnMint indicates if the signature is required on mint
func (s *SBTs) RequireSignatureOnMint() bool {
	return len(s.AuthorityPubKey) > 0
}

// Marshal marshals the SBTs
func (s *SBTs) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

// Unmarshal unmarshals the given data to the SBTs struct
func (s *SBTs) Unmarshal(data []byte) error {
	return json.Unmarshal(data, s)
}

// SBT defines the struct of the single token
type SBT struct {
	Symbol   string `json:"symbol"`             // unique symbol
	Id       uint64 `json:"token_id"`           // token id
	Owner    string `json:"owner"`              // token owner
	Metadata string `json:"metadata,omitempty"` // metadata per token

	BlockHeight         int64  `json:"block_height"` // mint block height
	TransactionIndex    int    `json:"tx_index"`     // mint tx index
	MintTransactionHash string `json:"mint_tx"`      // mint tx hash
}

// NewSBT creates a new SBT instance
func NewSBT(symbol string, id uint64, owner string, metadata string, blockHeight int64, txIndex int, mintTxHash string) *SBT {
	return &SBT{
		Symbol:              symbol,
		Id:                  id,
		Owner:               owner,
		Metadata:            metadata,
		BlockHeight:         blockHeight,
		TransactionIndex:    txIndex,
		MintTransactionHash: mintTxHash,
	}
}

// Compact returns the compact SBT
func (s *SBT) Compact() *CompactSBT {
	return &CompactSBT{
		Symbol:              s.Symbol,
		Id:                  s.Id,
		Metadata:            s.Metadata,
		MintTransactionHash: s.MintTransactionHash,
	}
}

// Marshal marshals the SBT
func (s *SBT) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

// Unmarshal unmarshals the given data to the SBT struct
func (s *SBT) Unmarshal(data []byte) error {
	return json.Unmarshal(data, s)
}

// CompactSBT defines the struct of the compact SBT token
type CompactSBT struct {
	Symbol              string `json:"symbol"`             // unique symbol
	Id                  uint64 `json:"token_id"`           // token id
	Metadata            string `json:"metadata,omitempty"` // metadata per token
	MintTransactionHash string `json:"mint_tx"`            // mint tx hash
}

// NewCompactSBT creates a new CompactSBT instance
func NewCompactSBT(symbol string, id uint64, metadata string, mintTxHash string) *CompactSBT {
	return &CompactSBT{
		Symbol:              symbol,
		Id:                  id,
		Metadata:            metadata,
		MintTransactionHash: mintTxHash,
	}
}

// Marshal marshals the CompactSBT
func (s *CompactSBT) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

// Unmarshal unmarshals the given data to the CompactSBT struct
func (s *CompactSBT) Unmarshal(data []byte) error {
	return json.Unmarshal(data, s)
}
