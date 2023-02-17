package commands

import (
	"fetch-go/pkg/fetch"

	"github.com/spf13/cobra"
)

var version = "0.0.1"

var rootCmd = &cobra.Command{
	Use:     "fetch",
	Version: version,
	Short:   "Fetch and download web pages",
	Long:    "A simple CLI to fetch and save web pages",
	Run: func(cmd *cobra.Command, args []string) {
		meta, _ := cmd.Flags().GetBool("metadata")
		fetch.Scrape(args, meta)
	},
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.Flags().BoolP("metadata", "m", false, "print metadata")
}
