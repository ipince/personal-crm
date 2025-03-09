package main

import (
	"fmt"

	"easyrm/pkg/cleanup"
	"easyrm/pkg/people"

	"github.com/spf13/cobra"
)

func setupEnrichCommands(rootCmd *cobra.Command) {
	enrichCmd := &cobra.Command{
		Use:   "enrich",
		Short: "Enrich contact data",
		Long:  "Enrich contacts with additional data (tags, gender, or all)",
	}
	rootCmd.AddCommand(enrichCmd)

	enrichTagsCmd := &cobra.Command{
		Use:   "tags",
		Short: "Apply tags to contacts",
		Run: func(cmd *cobra.Command, args []string) {
			srv := initializeClient()
			all := fetchAll(srv)
			cleanup.ApplyTags(srv, all)
			fmt.Println("Tags applied successfully")
		},
	}
	enrichCmd.AddCommand(enrichTagsCmd)

	enrichGenderCmd := &cobra.Command{
		Use:   "gender",
		Short: "Assign gender to contacts",
		Run: func(cmd *cobra.Command, args []string) {
			srv := initializeClient()
			all := fetchAll(srv)
			people.AssignGender(srv, all)
			fmt.Println("Gender assigned successfully")
		},
	}
	enrichCmd.AddCommand(enrichGenderCmd)

	// Enrich all subcommand
	enrichAllCmd := &cobra.Command{
		Use:   "all",
		Short: "Apply all enrichments to contacts",
		Run: func(cmd *cobra.Command, args []string) {
			srv := initializeClient()
			all := fetchAll(srv)

			// Apply tags
			cleanup.ApplyTags(srv, all)
			fmt.Println("Tags applied successfully")

			// Assign gender
			people.AssignGender(srv, all)
			fmt.Println("Gender assigned successfully")
		},
	}
	enrichCmd.AddCommand(enrichAllCmd)
}
