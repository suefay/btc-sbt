package params

// Params defines the BTC-SBT protocol params
type Params struct {
	ActivationBlockHeight int64 `json:"activationBlockHeight"` // activation height for the protocol
}

// NewParams creates a new Params instance
func NewParams(activationHeight int64) *Params {
	return &Params{
		ActivationBlockHeight: activationHeight,
	}
}

var (
	// MainNetParams is the params for the mainnet network
	MainNetParams = Params{
		ActivationBlockHeight: 823000,
	}

	// TestNetParams is the params for the testnet network
	TestNetParams = Params{
		ActivationBlockHeight: 2568000,
	}

	// SignetParams is the params for the signet network
	SigNetParams = Params{
		ActivationBlockHeight: 175000,
	}
)
