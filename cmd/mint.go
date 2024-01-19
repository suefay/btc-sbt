package cmd

import (
	"github.com/spf13/cobra"

	cfg "btc-sbt/config"
	"btc-sbt/initiator"
	"btc-sbt/protocol"
	"btc-sbt/stacks/basics"
)

func GetMintCmd() *cobra.Command {
	var addrType uint8

	cmd := &cobra.Command{
		Use:     "mint <symbol> <auth sig> <metadata> [flags] [config-file]",
		Short:   "Mint BTC SBT",
		Example: `btc-sbt mint sbt 0x123456 '{"attributes":[{"trait_type":"Level","value":"1"}]}'`,
		Args:    cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			configFileName := ""

			if len(args) == 3 {
				configFileName = cfg.DefaultConfigFileName
			} else {
				configFileName = args[3]
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

			op := protocol.NewMintOperation(args[0], addr.EncodeAddress(), args[1], args[2])

			commitTxHash, revealTxHash, err := initiator.Initiate(key, addr, op)
			if err != nil {
				return err
			}

			initiator.Logger.Infof("Minting SBT completed, commit tx: %s, reveal tx: %s", commitTxHash, revealTxHash)

			return nil
		},
	}

	cmd.Flags().Uint8VarP(&addrType, "addr-type", "a", 0, "address type; 0: taproot, 1: p2wpkh; default to taproot")

	return cmd
}
