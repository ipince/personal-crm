package main

import (
	"easyrm/people"
	"fmt"

	"github.com/spf13/cobra"
	peoplev1 "google.golang.org/api/people/v1"
)

// setupValidateCommands adds the validate commands to the root command
func setupValidateCommands(rootCmd *cobra.Command) {
	// Validate command
	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate contacts",
		Long:  "Validate different aspects of contacts (names, birthdays, or all)",
	}
	rootCmd.AddCommand(validateCmd)

	// Validate names subcommand
	validateNamesCmd := &cobra.Command{
		Use:   "names",
		Short: "Validate contact names",
		Run: func(cmd *cobra.Command, args []string) {
			srv := initializeClient()
			all := fetchAll(srv)
			validateNames(srv, all)
			fmt.Println("Name validation completed")
		},
	}
	validateCmd.AddCommand(validateNamesCmd)

	// Validate birthdays subcommand
	validateBirthdaysCmd := &cobra.Command{
		Use:   "birthdays",
		Short: "Validate contact birthdays",
		Run: func(cmd *cobra.Command, args []string) {
			srv := initializeClient()
			all := fetchAll(srv)
			validateBirthdays(srv, all)
			fmt.Println("Birthday validation completed")
		},
	}
	validateCmd.AddCommand(validateBirthdaysCmd)

	// Validate all subcommand
	validateAllCmd := &cobra.Command{
		Use:   "all",
		Short: "Validate all aspects of contacts",
		Run: func(cmd *cobra.Command, args []string) {
			srv := initializeClient()
			validateAll(srv)
			fmt.Println("Full validation completed")
		},
	}
	validateCmd.AddCommand(validateAllCmd)
}

// validateNames validates contact names
func validateNames(srv *peoplev1.Service, all []*peoplev1.Person) {
	for _, p := range all {
		err := people.ValidateNames(p)
		if err != nil {
			fmt.Printf("Name validation error for %s: %v\n", people.Name(p), err)
		}
	}
}

// validateBirthdays validates contact birthdays
func validateBirthdays(srv *peoplev1.Service, all []*peoplev1.Person) {
	for _, p := range all {
		err := people.ValidateBirthdays(p)
		if err != nil {
			fmt.Printf("Birthday validation error for %s: %v\n", people.Name(p), err)
		}
	}
}

// validateAll validates all aspects of contacts
func validateAll(srv *peoplev1.Service) {
	all := fetchAll(srv)
	for _, p := range all {
		people.Validate(p)
	}
}
