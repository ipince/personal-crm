package main

import (
	"fmt"
	"log"
	"os"

	"easyrm/pkg/contact"

	"github.com/spf13/cobra"
	peoplev1 "google.golang.org/api/people/v1"
)

var exportOutput string

// setupExportCommand adds the export command to the root command
func setupExportCommand(rootCmd *cobra.Command) {
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
