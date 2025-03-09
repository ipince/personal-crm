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
		Long:  "Summarize different aspects of contacts (all, genders, birthdays, addresses, urls, notes, tags)",
	}
	rootCmd.AddCommand(summarizeCmd)

	// Summarize all subcommand
	summarizeAllCmd := &cobra.Command{
		Use:   "all",
		Short: "Summarize all contact data",
		Run: func(cmd *cobra.Command, args []string) {
			srv := initializeClient()
			all := fetchAll(srv)
			err := cleanup.Summarize(srv, all)
			if err != nil {
				log.Fatalf("Failed to summarize: %v", err)
			}
		},
	}
	summarizeCmd.AddCommand(summarizeAllCmd)

	// Summarize genders subcommand
	summarizeGendersCmd := &cobra.Command{
		Use:   "genders",
		Short: "Summarize contact genders",
		Run: func(cmd *cobra.Command, args []string) {
			srv := initializeClient()
			all := fetchAll(srv)
			cleanup.SummarizeGenders(all)
		},
	}
	summarizeCmd.AddCommand(summarizeGendersCmd)

	// Summarize birthdays subcommand
	summarizeBirthdaysCmd := &cobra.Command{
		Use:   "birthdays",
		Short: "Summarize contact birthdays",
		Run: func(cmd *cobra.Command, args []string) {
			srv := initializeClient()
			all := fetchAll(srv)
			cleanup.SummarizeBirthdays(all)
		},
	}
	summarizeCmd.AddCommand(summarizeBirthdaysCmd)

	// Summarize addresses subcommand
	summarizeAddressesCmd := &cobra.Command{
		Use:   "addresses",
		Short: "Summarize contact addresses",
		Run: func(cmd *cobra.Command, args []string) {
			srv := initializeClient()
			all := fetchAll(srv)
			cleanup.SummarizeAddresses(all)
		},
	}
	summarizeCmd.AddCommand(summarizeAddressesCmd)

	// Summarize URLs subcommand
	summarizeUrlsCmd := &cobra.Command{
		Use:   "urls",
		Short: "Summarize contact URLs",
		Run: func(cmd *cobra.Command, args []string) {
			srv := initializeClient()
			all := fetchAll(srv)
			cleanup.SummarizeUrls(all)
		},
	}
	summarizeCmd.AddCommand(summarizeUrlsCmd)

	// Summarize notes subcommand
	summarizeNotesCmd := &cobra.Command{
		Use:   "notes",
		Short: "Summarize contact notes",
		Run: func(cmd *cobra.Command, args []string) {
			srv := initializeClient()
			all := fetchAll(srv)
			cleanup.SummarizeNotes(all)
		},
	}
	summarizeCmd.AddCommand(summarizeNotesCmd)

	// Summarize tags subcommand
	summarizeTagsCmd := &cobra.Command{
		Use:   "tags",
		Short: "Summarize contact tags",
		Run: func(cmd *cobra.Command, args []string) {
			srv := initializeClient()
			all := fetchAll(srv)
			err := cleanup.SummarizeTags(srv, all)
			if err != nil {
				log.Fatalf("Failed to summarize tags: %v", err)
			}
		},
	}
	summarizeCmd.AddCommand(summarizeTagsCmd)
}
