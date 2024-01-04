package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"btc-sbt/version"
)

func GetVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "BTC-SBT node and protocol version",
		Example: `btc-sbt version`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("node version: %s\nprotocol version: %s\n", version.NODE_VERSION, version.PROTOCOL_VERSION)

			return nil
		},
	}

	return cmd
}
