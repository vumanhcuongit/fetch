package commands

import (
	"fetch-go/pkg/fetch"
	"log"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var version = "0.0.1"

var rootCmd = &cobra.Command{
	Use:     "fetch",
	Version: version,
	Short:   "Fetch and download web pages",
	Long:    "A simple CLI to fetch and download web pages",
	Run: func(cmd *cobra.Command, args []string) {
		shouldFetchMetadata, _ := cmd.Flags().GetBool("metadata")
		maxConcurrent, _ := cmd.Flags().GetInt("max_concurrent")

		zaplogger, err := zap.NewProduction()
		if err != nil {
			log.Fatalf("failed to create logger: %v", err)
		}
		defer zaplogger.Sync()
		logger := zaplogger.Sugar()
		logger = logger.WithOptions(zap.AddCaller())

		fetcher := fetch.NewFetcher(logger)
		options := fetch.FetchOptions{ShouldFetchMetadata: shouldFetchMetadata, MaxConcurrent: maxConcurrent}
		fetcher.ScrapeURLs(args, options)
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
	rootCmd.Flags().Int("max_concurrent", 100, "max concurrent requests")
}
