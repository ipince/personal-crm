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

	"golang.org/x/oauth2/google"
	peoplev1 "google.golang.org/api/people/v1"
)

func main() {
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
	srv := people.Client(ctx, config)

	all := fetchAll(srv)
	err = cleanup.Summarize(srv, all)
	if err != nil {
		panic(err)
	}

	err = exportCsv(all, "export.csv")
	if err != nil {
		panic(err)
	}

	// Change things
	//err = addBirthdays(srv, all)
	//if err != nil {
	//	panic(err)
	//}

	//normalizeUrls(srv, all)
	//validateAll(srv)
	//people.AssignGender(srv, all)
	//cleanup.ApplyTags(srv, all)

	//validateAndNormalizeAll(srv)

}

func exportCsv(all []*peoplev1.Person, out string) error {
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, p := range all {
		c := contact.New(p)
		//notes := strings.ReplaceAll(c.Notes(), "\n", "")
		_, err = f.WriteString(fmt.Sprintf("%s,\"%s\",%s\n", c.FirstAndLastName(), c.FirstEmail(), c.FirstPhone()))
		if err != nil {
			return err
		}
	}
	return nil
}

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
		panic(err)
	}
}

func fetchAll(srv *peoplev1.Service) []*peoplev1.Person {
	all, err := people.ListAll(srv)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Fetched %d connections\n", len(all))
	return all
}

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

func validateAll(srv *peoplev1.Service) {
	all := fetchAll(srv)
	for _, p := range all {
		people.Validate(p)
	}
}

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
