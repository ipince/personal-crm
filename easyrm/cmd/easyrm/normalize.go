package main

import (
	"easyrm/people"
	"fmt"

	"github.com/spf13/cobra"
	peoplev1 "google.golang.org/api/people/v1"
)

// setupNormalizeCommands adds the normalize commands to the root command
func setupNormalizeCommands(rootCmd *cobra.Command) {
	// Normalize command
	normalizeCmd := &cobra.Command{
		Use:   "normalize",
		Short: "Normalize contact data",
		Long:  "Normalize different aspects of contacts (urls, phones, or all)",
	}
	rootCmd.AddCommand(normalizeCmd)

	// Normalize URLs subcommand
	normalizeUrlsCmd := &cobra.Command{
		Use:   "urls",
		Short: "Normalize contact URLs",
		Run: func(cmd *cobra.Command, args []string) {
			srv := initializeClient()
			all := fetchAll(srv)
			normalizeAllUrls(srv, all)
			fmt.Println("URLs normalized successfully")
		},
	}
	normalizeCmd.AddCommand(normalizeUrlsCmd)

	// Normalize phones subcommand
	normalizePhonesCmd := &cobra.Command{
		Use:   "phones",
		Short: "Normalize contact phone numbers",
		Run: func(cmd *cobra.Command, args []string) {
			srv := initializeClient()
			all := fetchAll(srv)
			normalizeAllPhones(srv, all)
			fmt.Println("Phone numbers normalized successfully")
		},
	}
	normalizeCmd.AddCommand(normalizePhonesCmd)

	// Normalize all subcommand
	normalizeAllCmd := &cobra.Command{
		Use:   "all",
		Short: "Normalize all aspects of contacts and validate",
		Run: func(cmd *cobra.Command, args []string) {
			srv := initializeClient()
			all := fetchAll(srv)
			normalizeAll(srv, all)
			fmt.Println("All normalizations completed")
		},
	}
	normalizeCmd.AddCommand(normalizeAllCmd)
}

// normalizeAllUrls normalizes URLs for all contacts using the people package implementation
func normalizeAllUrls(srv *peoplev1.Service, all []*peoplev1.Person) {
	for _, p := range all {
		err := people.NormalizeUrls(srv, p)
		if err != nil {
			fmt.Printf("Error normalizing URLs for %s: %v\n", people.Name(p), err)
		}
	}
}

// normalizeAllPhones normalizes phone numbers for all contacts using the existing implementation
func normalizeAllPhones(srv *peoplev1.Service, all []*peoplev1.Person) {
	for _, p := range all {
		err := people.Normalize(srv, p)
		if err != nil {
			fmt.Printf("Error normalizing phone numbers for %s: %v\n", people.Name(p), err)
		}
	}
}

// normalizeAll normalizes all aspects of contacts and validates them
func normalizeAll(srv *peoplev1.Service, all []*peoplev1.Person) {
	for _, p := range all {
		// Normalize URLs
		err := people.NormalizeUrls(srv, p)
		if err != nil {
			fmt.Printf("Error normalizing URLs for %s: %v\n", people.Name(p), err)
		}

		// Normalize phone numbers
		err = people.Normalize(srv, p)
		if err != nil {
			fmt.Printf("Error normalizing phone numbers for %s: %v\n", people.Name(p), err)
		}

		// Validate
		people.Validate(p)
	}
}
