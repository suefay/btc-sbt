package cmd

import (
	"encoding/hex"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/btcsuite/btcd/btcec/v2/schnorr"

	cfg "btc-sbt/config"
	"btc-sbt/initiator"
	"btc-sbt/protocol"
	"btc-sbt/stacks/basics"
)

func GetIssueCmd() *cobra.Command {
	var addrType uint8
	var selfPK bool

	cmd := &cobra.Command{
		Use:     "issue <symbol> <max supply> <auth pk> <end block> <metadata> [flags] [config-file]",
		Short:   "Issue BTC SBT",
		Example: `btc-sbt issue sbt 1000 0xabcdef 10000 '{"name":"sbt","description":""}‘`,
		Args:    cobra.RangeArgs(5, 6),
		RunE: func(cmd *cobra.Command, args []string) error {
			configFileName := ""

			if len(args) == 5 {
				configFileName = cfg.DefaultConfigFileName
			} else {
				configFileName = args[5]
			}

			maxSupply, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			endBlockHeight, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return err
			}

			v, err := cfg.LoadYAMLConfig(configFileName)
			if err != nil {
				return err
			}

			config, err := cfg.NewConfigFromViper(v)
			if err != nil {
				return err
			}

			initiator, err := initiator.NewInitiator(config)
			if err != nil {
				return err
			}

			key, addr, err := GetPrivateKeyAndAddress(config.KeyStorePath, basics.AddressType(addrType), initiator.NetParams)
			if err != nil {
				return err
			}

			authPK := args[2]
			if selfPK {
				authPK = hex.EncodeToString(schnorr.SerializePubKey(key.PubKey()))
			}

			op := protocol.NewIssueOperation(args[0], uint64(maxSupply), authPK, endBlockHeight, args[4])

			commitTxHash, revealTxHash, err := initiator.Initiate(key, addr, op)
			if err != nil {
				return err
			}

			initiator.Logger.Infof("Issuing SBT completed, commit tx: %s, reveal tx: %s", commitTxHash, revealTxHash)

			return nil
		},
	}

	cmd.Flags().Uint8VarP(&addrType, "addr-type", "a", 0, "address type; 0: taproot, 1: p2wpkh; default to taproot")
	cmd.Flags().BoolVarP(&selfPK, "self-pk", "p", false, "indicates if the current public key is used for verification")

	return cmd
}
