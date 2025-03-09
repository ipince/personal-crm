package main

import (
	"log"

	"easyrm/pkg/cleanup"

	"github.com/spf13/cobra"
)

// setupSummarizeCommand adds the summarize command to the root command
func setupSummarizeCommand(rootCmd *cobra.Command) {
	// Summarize command
	summarizeCmd := &cobra.Command{
		Use:   "summarize",
		Short: "Summarize contact data",
		Run: func(cmd *cobra.Command, args []string) {
			srv := initializeClient()
			all := fetchAll(srv)
			err := cleanup.Summarize(srv, all)
			if err != nil {
				log.Fatalf("Failed to summarize: %v", err)
			}
		},
	}
	rootCmd.AddCommand(summarizeCmd)
}
