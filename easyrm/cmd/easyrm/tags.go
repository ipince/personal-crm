package main

import (
	"fmt"

	"easyrm/pkg/cleanup"
	"easyrm/pkg/people"

	"github.com/spf13/cobra"
)

// setupTagsCommands adds the tags and gender commands to the root command
func setupTagsCommands(rootCmd *cobra.Command) {
	// Apply tags command
	applyTagsCmd := &cobra.Command{
		Use:   "apply-tags",
		Short: "Apply tags to contacts",
		Run: func(cmd *cobra.Command, args []string) {
			srv := initializeClient()
			all := fetchAll(srv)
			cleanup.ApplyTags(srv, all)
			fmt.Println("Tags applied successfully")
		},
	}
	rootCmd.AddCommand(applyTagsCmd)

	// Assign gender command
	assignGenderCmd := &cobra.Command{
		Use:   "assign-gender",
		Short: "Assign gender to contacts",
		Run: func(cmd *cobra.Command, args []string) {
			srv := initializeClient()
			all := fetchAll(srv)
			people.AssignGender(srv, all)
			fmt.Println("Gender assigned successfully")
		},
	}
	rootCmd.AddCommand(assignGenderCmd)
}
