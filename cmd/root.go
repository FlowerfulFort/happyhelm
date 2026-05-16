package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "happyhelm <chart> [keyword...]",
	Short: "Interactive Helm values picker",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPick(args[0], args[1:], pickOpts)
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	addPickFlags(rootCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
