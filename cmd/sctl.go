package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var (
	inputFilePath string
	sctlCmd       = &cobra.Command{
		Use:   "sctl",
		Short: "sctl is a Selenium configuration management tool",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}
)

func Execute() {
	sctlCmd.AddCommand(generateCmd)
	sctlCmd.AddCommand(statCmd)

	if _, err := sctlCmd.ExecuteC(); err != nil {
		os.Exit(1)
	}
}

func initCommonFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&inputFilePath, "inputFile", "input.json", "path to input file")
}
