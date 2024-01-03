package cmd

import (
	"github.com/spf13/cobra"

	cfg "btc-sbt/config"
	"btc-sbt/node"
)

func GetNodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "node [config-file]",
		Short:   "Start BTC-SBT node",
		Example: `btc-sbt node`,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			configFileName := ""

			if len(args) == 0 {
				configFileName = cfg.DefaultConfigFileName
			} else {
				configFileName = args[0]
			}

			v, err := cfg.LoadYAMLConfig(configFileName)
			if err != nil {
				return err
			}

			config, err := cfg.NewConfigFromViper(v)
			if err != nil {
				return err
			}

			node, err := node.CreateNode(config)
			if err != nil {
				return err
			}

			node.Start()

			return nil
		},
	}

	return cmd
}
