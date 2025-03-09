package main

import (
	"context"
	"easyrm/people"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/oauth2/google"
	peoplev1 "google.golang.org/api/people/v1"
)

var rootCmd = &cobra.Command{
	Use:   "easyrm",
	Short: "EasyRM - Personal CRM CLI",
	Long:  `EasyRM is a command-line tool for managing your Google Contacts as a personal CRM.`,
}

func init() {
	setupExportCommand(rootCmd)
	setupNormalizeCommands(rootCmd)
	setupValidateCommands(rootCmd)
	setupSummarizeCommand(rootCmd)
	setupListCommands(rootCmd)
	setupTagsCommands(rootCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// initializeClient initializes the Google People API client
func initializeClient() *peoplev1.Service {
	ctx := context.Background()
	b, err := os.ReadFile("oauth_desktop_client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, peoplev1.ContactsScope) // ContactsReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	return people.Client(ctx, config)
}
