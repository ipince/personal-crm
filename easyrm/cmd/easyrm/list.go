package main

import (
	"fmt"
	"log"
	"strings"

	"easyrm/pkg/people"

	"github.com/spf13/cobra"
	peoplev1 "google.golang.org/api/people/v1"
)

// setupListCommands adds list commands to the root command
func setupListCommands(rootCmd *cobra.Command) {
	// List command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List contacts",
		Long:  "List contacts with different filters (all, z-names)",
	}
	rootCmd.AddCommand(listCmd)

	// List all subcommand
	listAllCmd := &cobra.Command{
		Use:   "all",
		Short: "List all contacts",
		Run: func(cmd *cobra.Command, args []string) {
			srv := initializeClient()
			fetchAll(srv)
		},
	}
	listCmd.AddCommand(listAllCmd)

	// List Z names subcommand
	listZNamesCmd := &cobra.Command{
		Use:   "z-names",
		Short: "List contacts with 'z' in their name",
		Run: func(cmd *cobra.Command, args []string) {
			srv := initializeClient()
			all := fetchAll(srv)
			listZNames(all)
		},
	}
	listCmd.AddCommand(listZNamesCmd)

	// Remove the old list-z-names command from the root command
	// This is now a subcommand of the list command
}

// fetchAll fetches all contacts
func fetchAll(srv *peoplev1.Service) []*peoplev1.Person {
	all, err := people.ListAll(srv)
	if err != nil {
		log.Fatalf("Failed to fetch contacts: %v", err)
	}
	fmt.Printf("Fetched %d contacts\n", len(all))
	return all
}

// listZNames lists contacts with 'z' in their name
func listZNames(peeps []*peoplev1.Person) {
	count := 0
	for _, p := range peeps {
		if len(p.Names) > 0 {
			parts := strings.Split(p.Names[0].DisplayName, " ")
			for _, s := range parts {
				if strings.ToLower(s) == "z" {
					count++
					fmt.Printf("found person with z name: %s (%s)\n", p.Names[0].DisplayName, people.Link(p))
				}
			}
		}
	}
	fmt.Printf("found %d names\n", count)
}
