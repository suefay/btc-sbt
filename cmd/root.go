package cmd

import (
	"github.com/spf13/cobra"
)

func GetRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "btc-sbt",
		Short:   "Command Line Interface for the BTC-SBT protocol",
		Example: `btc-sbt -h`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	return cmd
}
