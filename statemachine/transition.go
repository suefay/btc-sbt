package statemachine

import (
	"encoding/hex"
	"fmt"

	"btc-sbt/crypto/signature/schnorr"
	"btc-sbt/protocol"
	"btc-sbt/types"
)

// HandleOps handles the specified protocol operations in the given context
func (sm *StateMachine) HandleOps(ctx *Context, ops []protocol.Operation) error {
	for _, op := range ops {
		err := sm.HandleOp(ctx, op)
		if err != nil && IsExecutionFailedErr(err) {
			return err
		}
	}

	return nil
}

// HandleOp handles the specified protocol operation in the given context
func (sm *StateMachine) HandleOp(ctx *Context, op protocol.Operation) error {
	switch op := op.(type) {
	case *protocol.IssueOperation:
		return sm.HandleIssue(ctx, op)

	case *protocol.MintOperation:
		return sm.HandleMint(ctx, op)
	}

	return nil
}

// HandleIssue handles the state transition for the issue operation
func (sm *StateMachine) HandleIssue(ctx *Context, op *protocol.IssueOperation) error {
	if err := op.Validate(sm.NetParams); err != nil {
		return wrapError(InvalidOpErr, err)
	}

	ok, err := sm.SBTsExists(op.Symbol)
	if err != nil {
		return wrapError(ExecutionFailedErr, err)
	}

	if ok {
		return wrapError(InvalidOpErr, fmt.Errorf("symbol already exists: %s", op.Symbol))
	}

	seq, err := sm.IncreaseSBTsSequence()
	if err != nil {
		return wrapError(ExecutionFailedErr, err)
	}

	sbts := types.NewSBTsFromIssueOp(op)

	sbts.Sequence = seq
	sbts.Issuer = ctx.OperationOutAddress
	sbts.BlockHeight = ctx.BlockHeight
	sbts.TransactionIndex = ctx.TxIndex
	sbts.IssueTransactionHash = ctx.Tx.TxHash().String()

	if err := sm.SetSBTs(sbts); err != nil {
		return wrapError(ExecutionFailedErr, err)
	}

	return nil
}

// HandleMint handles the state transition for the mint operation
func (sm *StateMachine) HandleMint(ctx *Context, op *protocol.MintOperation) error {
	if err := op.Validate(sm.NetParams); err != nil {
		return wrapError(InvalidOpErr, err)
	}

	sbts, err := sm.GetSBTs(op.Symbol)
	if err != nil {
		return wrapError(ExecutionFailedErr, err)
	}

	if sbts == nil {
		return wrapError(InvalidOpErr, fmt.Errorf("symbol does not exist: %s", op.Symbol))
	}

	if sbts.EndBlockHeight > 0 && ctx.BlockHeight > sbts.EndBlockHeight {
		return wrapError(InvalidOpErr, fmt.Errorf("mint ended at block %d", sbts.EndBlockHeight))
	}

	if sbts.MaxSupply > 0 && sbts.TotalSupply+1 > sbts.MaxSupply {
		return wrapError(InvalidOpErr, fmt.Errorf("max supply reached: %d", sbts.MaxSupply))
	}

	ok, err := sm.HasOwnedSBT(op.Owner, op.Symbol)
	if err != nil {
		return wrapError(ExecutionFailedErr, err)
	}

	if ok {
		return wrapError(InvalidOpErr, fmt.Errorf("address has owned the SBT: %s, %s", op.Owner, op.Symbol))
	}

	if sbts.RequireSignatureOnMint() {
		if len(op.AuthoritySignature) == 0 {
			return wrapError(InvalidOpErr, fmt.Errorf("authority signature required"))
		}

		sigHash, err := op.Hash()
		if err != nil {
			return wrapError(ExecutionFailedErr, err)
		}

		// validated
		sigBytes, _ := hex.DecodeString(op.AuthoritySignature)
		pubKeyBytes, _ := hex.DecodeString(sbts.AuthorityPubKey)

		if !schnorr.VerifySignature(sigBytes, sigHash, pubKeyBytes) {
			return wrapError(InvalidOpErr, fmt.Errorf("authority signature verification failed"))
		}
	}

	sbt := types.NewSBT(op.Symbol, sbts.TotalSupply, op.Owner, op.Metadata, ctx.BlockHeight, ctx.TxIndex, ctx.Tx.TxHash().String())

	if err := sm.SetSBT(sbt); err != nil {
		return wrapError(ExecutionFailedErr, err)
	}

	if err := sm.SetOwnerSBT(op.Owner, sbt); err != nil {
		return wrapError(ExecutionFailedErr, err)
	}

	if err := sm.SetSBTsSupply(op.Symbol, sbts.TotalSupply+1); err != nil {
		return wrapError(ExecutionFailedErr, err)
	}

	return nil
}
