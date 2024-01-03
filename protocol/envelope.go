package protocol

import (
	"bytes"

	"github.com/btcsuite/btcd/txscript"

	"btc-sbt/stacks/basics"
)

// Envelope is the envelope of the BTC-SBT protocol
type Envelope struct {
	Header             []byte // OP_FALSE OP_IF
	ProtocolIdentifier []byte // protocol identifier
	PayloadTag         byte   // OP_FALSE
	Payload            []byte // protocol payload
	Tail               []byte // OP_ENDIF
}

// NewEnvelope builds a new Envelope instance with the specified payload
func NewEnvelope(payload []byte) *Envelope {
	return &Envelope{
		Header:             ENVELOPE_HEADER,
		ProtocolIdentifier: ENVELOPE_PROTOCOL_IDENTIFIER,
		PayloadTag:         ENVELOPE_PAYLOAD_TAG,
		Payload:            payload,
		Tail:               ENVELOPE_TAIL,
	}
}

// Script gets the script representation of the envelope
func (e *Envelope) Script() ([]byte, error) {
	scriptBuilder := txscript.NewScriptBuilder()

	scriptBuilder.AddOps(e.Header)
	scriptBuilder.AddData(e.ProtocolIdentifier)
	scriptBuilder.AddOp(e.PayloadTag)

	basics.AddLargeDataToScript(scriptBuilder, e.Payload)

	scriptBuilder.AddOps(e.Tail)

	script, err := scriptBuilder.Script()
	if err != nil {
		return nil, err
	}

	return script, nil
}

// ExtractEnvelope extracts the BTC-SBT protocol envelope from the given script
func ExtractEnvelope(script []byte) *Envelope {
	tokenizer := txscript.MakeScriptTokenizer(0, script)

	hasHeader := false
	hasProtocolId := false
	hasPayloadTag := false
	hasTail := false

	payload := make([]byte, 0)

	for tokenizer.Next() {
		if tokenizer.Opcode() == ENVELOPE_HEADER[0] && tokenizer.Next() && tokenizer.Opcode() == ENVELOPE_HEADER[1] {
			hasHeader = true

			break
		}
	}

	if hasHeader && tokenizer.Next() && bytes.Equal(tokenizer.Data(), ENVELOPE_PROTOCOL_IDENTIFIER) {
		hasProtocolId = true
	}

	if hasProtocolId && tokenizer.Next() && tokenizer.Opcode() == ENVELOPE_PAYLOAD_TAG {
		hasPayloadTag = true
	}

	if hasPayloadTag {
		for tokenizer.Next() {
			payload = append(payload, tokenizer.Data()...)

			if tokenizer.Opcode() == ENVELOPE_TAIL[0] {
				hasTail = true

				break
			}
		}
	}

	if hasTail && len(payload) > 0 {
		return NewEnvelope(payload)
	}

	return nil
}
