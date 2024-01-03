package cmd

import (
	"strconv"

	"github.com/spf13/cobra"

	cfg "btc-sbt/config"
	"btc-sbt/initiator"
	"btc-sbt/protocol"
)

func GetIssueCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "issue <symbol> <max supply> <pk> <end block> <metadata> [config-file]",
		Short:   "Issue BTC SBT",
		Example: `btc-sbt issue sbt 1000 10000 0xabcdef '{"name":"sbt","description":""}‘`,
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

			key, _, err := GetPrivateKeyAndTaprootAddr(config.KeyStorePath, initiator.NetParams)
			if err != nil {
				return err
			}

			op := protocol.NewIssueOperation(args[0], uint64(maxSupply), args[2], endBlockHeight, args[4])

			commitTxHash, revealTxHash, err := initiator.Initiate(key, op)
			if err != nil {
				return err
			}

			initiator.Logger.Infof("Issuing SBT completed, commit tx: %s, reveal tx: %s", commitTxHash, revealTxHash)

			return nil
		},
	}

	return cmd
}
