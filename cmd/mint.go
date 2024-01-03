package cmd

import (
	"github.com/spf13/cobra"

	cfg "btc-sbt/config"
	"btc-sbt/initiator"
	"btc-sbt/protocol"
)

func GetMintCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mint <symbol> <authsig> <metadata> [config-file]",
		Short:   "Mint BTC SBT",
		Example: `btc-sbt mint sbt 0x123456 '{"level":1}' [config-file]`,
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

			key, taprootAddr, err := GetPrivateKeyAndTaprootAddr(config.KeyStorePath, initiator.NetParams)
			if err != nil {
				return err
			}

			op := protocol.NewMintOperation(args[0], taprootAddr.EncodeAddress(), args[1], args[2])

			commitTxHash, revealTxHash, err := initiator.Initiate(key, op)
			if err != nil {
				return err
			}

			initiator.Logger.Infof("Minting SBT completed, commit tx: %s, reveal tx: %s", commitTxHash, revealTxHash)

			return nil
		},
	}

	return cmd
}
