package people

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"

	peoplev1 "google.golang.org/api/people/v1"
)

type facebookContact struct {
	Name string
	URL  string
}

// MergeFacebookURLs reads a csv file of Facebook URLs and names, and attempts to merge them
// into a slice of People.
func MergeFacebookURLs(srv *peoplev1.Service, all []*peoplev1.Person) {

	// Load up the (fbname, fburl) pairs
	f, err := os.Open("data/fb_export_2021-06-30.csv") // TODO: run with new dataset
	if err != nil {
		panic(err)
	}
	records, err := csv.NewReader(f).ReadAll()
	fbFriends := []*facebookContact{}
	for _, r := range records {
		fbFriends = append(fbFriends, &facebookContact{
			Name: r[1],
			URL:  r[2],
		})
	}

	// Get contacts who already have an FB URL. These will not be edited. ... well, depends.
	fbContacts := map[string]*peoplev1.Person{}
	for _, p := range all {
		if u := facebookURL(p); u != nil {
			fbContacts[u.Value] = p
		}
	}
	fmt.Printf("Found %d contacts with fb links\n", len(fbContacts))

	// For each fbFriend, see if we can find a contact who matches
	for _, fbf := range fbFriends {
		// Skip matching for contacts who already have a Facebook URL.
		if _, ok := fbContacts[fbf.URL]; ok {
			continue
		}

		// Match by name. If found, add Facebook URL and be done.
		found := false
		for _, c := range all {
			if len(c.Names) == 0 { // must be org, ignore.
				continue
			}
			if fbf.Name == c.Names[0].DisplayName {
				// Easy match!
				fmt.Printf("easy match! %s -> %s\n", fbf.URL, Link(c))
				err = setFacebookURL(srv, c, fbf.URL)
				if err != nil {
					fmt.Println("ERR: " + err.Error())
				}
				found = true
				break
			}
		}
		if found {
			continue
		}

		// Not found by name, so query the Google People service.
		// - If there's a single match, adds the URL to the person (assuming it doesn't already have one).
		// - If there are no matches, a new person is created and tagged with "Staging". The user may want
		// to review these and merge them with another person. This case is common when the person uses
		// a fake name on Facebook.
		// TODO: remove the "Staging" tagging, since it's specific to my account!
		// - If there are multiple matches, we print that out and let the user deal with it manually. This case
		// may happen if a person has multiple fb profiles; or if two people actually have the same name.
		time.Sleep(1 * time.Second)
		r, err := srv.People.SearchContacts().Query(fbf.Name).ReadMask("names").Do()
		if err != nil {
			fmt.Println("failed search query: " + err.Error())
			continue
		}
		if len(r.Results) == 0 {
			fmt.Printf("no query match found for %s (%s). adding new person\n", fbf.Name, fbf.URL)
			// No match found, let's add it!
			err = Insert(srv, fbf.Name, fbf.URL)
			if err != nil {
				fmt.Println("ERR: " + err.Error())
				continue
			}
		} else if len(r.Results) > 1 {
			summary := []string{}
			for _, r := range r.Results {
				summary = append(summary, Link(r.Person))
			}
			fmt.Printf("too many results for %s (%s): %s\n", fbf.Name, fbf.URL, strings.Join(summary, ", "))
		} else {
			// exactly one match
			fmt.Printf("query match! %s -> %s\n", fbf.Name, Link(r.Results[0].Person))
			if facebookURL(r.Results[0].Person) != nil {
				err = setFacebookURL(srv, r.Results[0].Person, fbf.URL)
				if err != nil {
					fmt.Println("ERR: " + err.Error())
					continue
				}
			}
		}
	}
}
