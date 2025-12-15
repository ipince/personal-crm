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

	// List by location subcommand
	listLocationCmd := &cobra.Command{
		Use:   "location [location]",
		Short: "List contacts by location (city, state, or country)",
		Long:  "List contacts filtered by location. Location can be a city, US state, or country. Searches in the 'Current City' UserDefined field.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			srv := initializeClient()
			all := fetchAll(srv)
			listByLocation(all, args[0])
		},
	}
	listCmd.AddCommand(listLocationCmd)
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

// listByLocation lists contacts filtered by location
func listByLocation(peeps []*peoplev1.Person, location string) {
	count := 0
	location = strings.ToLower(location)
	
	for _, p := range peeps {
		currentCity := getCurrentCity(p)
		if currentCity != "" && matchesLocation(currentCity, location) {
			count++
			fmt.Printf("%s - %s (%s)\n", people.Name(p), currentCity, people.Link(p))
		}
	}
	fmt.Printf("found %d contacts in %s\n", count, location)
}

// getCurrentCity extracts the current city from a person's UserDefined fields
func getCurrentCity(p *peoplev1.Person) string {
	for _, ud := range p.UserDefined {
		if ud.Key == "Current City" {
			return ud.Value
		}
	}
	return ""
}

// matchesLocation checks if a location string matches the search term
// Supports matching city, state, or country
func matchesLocation(location, searchTerm string) bool {
	location = strings.ToLower(location)
	searchTerm = strings.ToLower(searchTerm)
	
	// Split location by comma to get city, state/country parts
	parts := strings.Split(location, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.Contains(part, searchTerm) {
			return true
		}
	}
	return false
}
