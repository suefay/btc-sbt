package main

import (
	"fmt"
	"os"

	"btc-sbt/cmd"
)

func main() {
	rootCmd := cmd.GetRootCmd()

	nodeCmd := cmd.GetNodeCmd()

	issueCmd := cmd.GetIssueCmd()
	mintCmd := cmd.GetMintCmd()

	versionCmd := cmd.GetVersionCmd()

	rootCmd.AddCommand(nodeCmd)
	rootCmd.AddCommand(issueCmd)
	rootCmd.AddCommand(mintCmd)
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
}
