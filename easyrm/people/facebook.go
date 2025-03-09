package people

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"

	peoplev1 "google.golang.org/api/people/v1"
)

type FacebookFriend struct {
	Name string
	URL  string
}

// MergeFacebookURLs reads a csv file of Facebook URLs and names, and attempts to merge them
// into a slice of People.
func MergeFacebookURLs(srv *peoplev1.Service, all []*peoplev1.Person, fbFriends []*FacebookFriend) {

	// Get contacts who already have an FB URL. These will not be edited.
	fbContacts := map[string]*peoplev1.Person{}
	for _, p := range all {
		for _, u := range FacebookURLs(p) {
			fbContacts[u] = p
		}
	}
	fmt.Printf("Found %d contacts with fb links\n", len(fbContacts))

	// For each fbFriend, see if we can find a contact who matches
	for _, fbf := range fbFriends {
		// Skip matching if we already have a contact with this Facebook URL.
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
				found = true
				// Easy match!
				fmt.Printf("easy match! %s -> %s\n", fbf.URL, Link(c))
				err := setFacebookURLIfNoneExists(srv, c, fbf.URL)
				if err != nil {
					fmt.Println("ERR: " + err.Error())
				}
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
		r, err := srv.People.SearchContacts().Query(fbf.Name).ReadMask("names,urls").Do()
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
		} else { // exactly one match
			fmt.Printf("query match! %s -> %s\n", fbf.URL, Link(r.Results[0].Person))
			err = setFacebookURLIfNoneExists(srv, r.Results[0].Person, fbf.URL)
			if err != nil {
				fmt.Println("ERR: " + err.Error())
				continue
			}
		}
	}
}

func LoadFacebookFriends(csvPath string) ([]*FacebookFriend, error) {
	// Load up the (fbname, fburl) pairs
	f, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}
	friends := []*FacebookFriend{}
	for _, r := range records {
		friends = append(friends, &FacebookFriend{
			Name: strings.TrimSpace(r[1]),
			URL:  strings.TrimSpace(r[0]),
		})
	}
	return friends, nil
}

// FacebookURLs returns a slice of urls in person that are Facebook links.
func FacebookURLs(person *peoplev1.Person) []string {
	urls := []string{}
	for _, u := range person.Urls {
		if strings.Contains(u.Value, "facebook.com") { // could also match on Type.
			urls = append(urls, u.Value)
		}
	}
	return urls
}

func setFacebookURLIfNoneExists(srv *peoplev1.Service, person *peoplev1.Person, url string) error {
	if len(FacebookURLs(person)) == 0 {
		err := setFacebookURL(srv, person, url)
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("skipping adding %s to %s because it would overwrite existing url\n", url, Link(person))
	}
	return nil
}

func setFacebookURL(srv *peoplev1.Service, person *peoplev1.Person, fbURL string) error {

	// Check if one exists and replace if so.
	for _, u := range person.Urls {
		if strings.Contains(u.Value, "facebook.com") {
			// replace
			u.Type = "Facebook" // maybe make it say facebook?
			u.Value = fbURL
			return update(srv, person, "urls")
		}
	}

	// else, add it
	person.Urls = append(person.Urls, &peoplev1.Url{
		Type:  "Facebook",
		Value: fbURL,
	})
	return update(srv, person, "urls")
}
