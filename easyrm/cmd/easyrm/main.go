package main

import (
	"context"
	"easyrm/cleanup"
	"easyrm/contact"
	"easyrm/people"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/oauth2/google"
	peoplev1 "google.golang.org/api/people/v1"
)

var (
	exportOutput string
	rootCmd      = &cobra.Command{
		Use:   "easyrm",
		Short: "EasyRM - Personal CRM CLI",
		Long:  `EasyRM is a command-line tool for managing your Google Contacts as a personal CRM.`,
	}
)

func init() {
	// Export command
	exportCmd := &cobra.Command{
		Use:   "export",
		Short: "Export contacts to CSV",
		Run: func(cmd *cobra.Command, args []string) {
			srv := initializeClient()
			all := fetchAll(srv)
			err := exportCsv(all, exportOutput)
			if err != nil {
				log.Fatalf("Failed to export CSV: %v", err)
			}
			fmt.Printf("Exported contacts to %s\n", exportOutput)
		},
	}
	exportCmd.Flags().StringVarP(&exportOutput, "out", "o", "export.csv", "Output CSV file path")
	rootCmd.AddCommand(exportCmd)

	// Normalize URLs command
	normalizeUrlsCmd := &cobra.Command{
		Use:   "normalize-urls",
		Short: "Normalize URLs in contacts",
		Run: func(cmd *cobra.Command, args []string) {
			srv := initializeClient()
			all := fetchAll(srv)
			normalizeUrls(srv, all)
			fmt.Println("URLs normalized successfully")
		},
	}
	rootCmd.AddCommand(normalizeUrlsCmd)

	// Fetch all command
	fetchAllCmd := &cobra.Command{
		Use:   "fetch-all",
		Short: "Fetch all contacts",
		Run: func(cmd *cobra.Command, args []string) {
			srv := initializeClient()
			all := fetchAll(srv)
			fmt.Printf("Fetched %d contacts\n", len(all))
		},
	}
	rootCmd.AddCommand(fetchAllCmd)

	// List Z names command
	listZNamesCmd := &cobra.Command{
		Use:   "list-z-names",
		Short: "List contacts with 'z' in their name",
		Run: func(cmd *cobra.Command, args []string) {
			srv := initializeClient()
			all := fetchAll(srv)
			listZNames(all)
		},
	}
	rootCmd.AddCommand(listZNamesCmd)

	// Validate command
	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate all contacts",
		Run: func(cmd *cobra.Command, args []string) {
			srv := initializeClient()
			validateAll(srv)
			fmt.Println("Validation completed")
		},
	}
	rootCmd.AddCommand(validateCmd)

	// Validate and normalize command
	validateAndNormalizeCmd := &cobra.Command{
		Use:   "validate-and-normalize",
		Short: "Validate and normalize all contacts",
		Run: func(cmd *cobra.Command, args []string) {
			srv := initializeClient()
			validateAndNormalizeAll(srv)
			fmt.Println("Validation and normalization completed")
		},
	}
	rootCmd.AddCommand(validateAndNormalizeCmd)

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

// exportCsv exports contacts to a CSV file
func exportCsv(all []*peoplev1.Person, out string) error {
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, p := range all {
		c := contact.New(p)
		_, err = f.WriteString(fmt.Sprintf("%s,\"%s\",%s\n", c.FirstAndLastName(), c.FirstEmail(), c.FirstPhone()))
		if err != nil {
			return err
		}
	}
	return nil
}

// normalizeUrls normalizes URLs in contacts
func normalizeUrls(srv *peoplev1.Service, all []*peoplev1.Person) {
	domains := map[string]string{
		"facebook.com":  "Facebook",
		"linkedin.com":  "LinkedIn",
		"instagram.com": "Instagram",
		"youtube.com":   "YouTube",
	}
	remove := "plus.google.com"
	save := []*peoplev1.Person{}
	for _, p := range all {
		rmIndex := -1
		dirty := false
		for i, u := range p.Urls {
			if u.Metadata.Source.Type != "CONTACT" {
				continue
			}

			if strings.Contains(u.Value, remove) {
				rmIndex = i
				dirty = true
				continue
			}

			known := false
			for domain, label := range domains {
				if strings.Contains(u.Value, domain) {
					known = true
					if u.Type != label {
						fmt.Printf("%s: %s (%s -> %s)\n", people.Name(p), u.Value, u.Type, label)
						u.Type = label
						dirty = true
						break
					}
				}
			}
			if !known && u.Type != "other" {
				fmt.Printf("%s: %s (%s -> %s)\n", people.Name(p), u.Value, u.Type, "other")
				u.Type = "other"
				dirty = true
			}
		}
		if rmIndex != -1 { // remove...
			p.Urls = append(p.Urls[:rmIndex], p.Urls[rmIndex+1:]...)
		}
		if dirty {
			save = append(save, p)
		}
	}
	err := people.UpdateAll(srv, save, "urls")
	if err != nil {
		log.Fatalf("Failed to update URLs: %v", err)
	}
}

// fetchAll fetches all contacts
func fetchAll(srv *peoplev1.Service) []*peoplev1.Person {
	all, err := people.ListAll(srv)
	if err != nil {
		log.Fatalf("Failed to fetch contacts: %v", err)
	}
	fmt.Printf("Fetched %d connections\n", len(all))
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

// validateAll validates all contacts
func validateAll(srv *peoplev1.Service) {
	all := fetchAll(srv)
	for _, p := range all {
		people.Validate(p)
	}
}

// validateAndNormalizeAll validates and normalizes all contacts
func validateAndNormalizeAll(srv *peoplev1.Service) {
	all := fetchAll(srv)
	for _, p := range all {
		err := people.Normalize(srv, p)
		if err != nil {
			fmt.Println(err)
		}
		people.Validate(p)
	}
}
