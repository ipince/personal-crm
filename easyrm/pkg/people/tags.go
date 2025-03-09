package people

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	peoplev1 "google.golang.org/api/people/v1"
)

type Tag struct {
	Name                string
	Description         string
	Options             []Option
	Exhaustive          bool // if true, all eligible members must choose an option
	CreateIfNonExistent bool
	Applicable          func(p *peoplev1.Person) bool
	Rules               []Rule
	// TODO: implement a MutuallyExclusive bool param
}
type Option struct {
	Name        string
	Description string
}

type Rule func(p *peoplev1.Person) (string, error)

func (t *Tag) AppliesTo(p *peoplev1.Person) bool {
	if t.Applicable != nil {
		return t.Applicable(p)
	}
	return true
}

func AskWithPrompt(tag *Tag) Rule {
	return func(p *peoplev1.Person) (string, error) {
		ans := ""
		optionDescs := []string{}
		optionNames := map[string]string{}
		for _, o := range tag.Options {
			optionDescs = append(optionDescs, o.Description)
			optionNames[o.Description] = o.Name
		}
		prompt := fmt.Sprintf("%s | %s | %s | %s | %s", Name(p), Updated(p), Note(p), Link(p), Facebook(p))
		err := survey.AskOne(&survey.Select{
			Message: prompt,
			Options: append(optionDescs, "skip"),
			Help:    fmt.Sprintf("Which option fits %s?", tag.Description),
		}, &ans)
		if err != nil {
			// user cancellation; pass error so we can exit -> need to pass err?
			return "", err
		}
		if ans == "skip" {
			ans = ""
		}
		return optionNames[ans], nil
	}
}

func ApplyTag(srv *peoplev1.Service, people []*peoplev1.Person, tag *Tag) error {
	// shuffle
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(people), func(i, j int) { people[i], people[j] = people[j], people[i] })

	// Get Group for each option.
	allGroups, err := srv.ContactGroups.List().PageSize(100).Do()
	if err != nil {
		return err
	}
	groupNames := map[string]string{}             // id -> formatted name
	groups := map[string]*peoplev1.ContactGroup{} // name -> group
	for _, g := range allGroups.ContactGroups {
		groupNames[g.ResourceName] = g.FormattedName
		for _, t := range tag.Options {
			if g.FormattedName == t.Name { // assume there are no duplicates
				groups[t.Name] = g
				break
			}
		}
	}

	// Create the ones that don't exist yet, if needed.
	for _, t := range tag.Options {
		if _, ok := groups[t.Name]; !ok {
			if tag.CreateIfNonExistent {
				fmt.Printf("Creating a new Group for option %s\n", t.Name)
				newGroup, err := srv.ContactGroups.Create(&peoplev1.CreateContactGroupRequest{
					ContactGroup: &peoplev1.ContactGroup{
						Name: t.Name,
					},
				}).Do()
				if err != nil {
					return err
				}
				// add to maps
				groupNames[newGroup.ResourceName] = newGroup.FormattedName
				groups[t.Name] = newGroup
			}
		}
	}

	// Go through people and try to put them in an option.
	batch := []*peoplev1.Person{}
	batchSize := 10
	for _, p := range people {
		// Skip ineligible contacts
		if !tag.AppliesTo(p) {
			continue
		}

		// Check if person is already a member of an option.
		alreadyTagged := false
		for _, g := range p.Memberships {
			for _, opt := range tag.Options {
				if groupNames[g.ContactGroupMembership.ContactGroupResourceName] == opt.Name {
					//fmt.Printf("%s is already tagged with option %s. skipping...\n", Name(p), opt.Name)
					alreadyTagged = true
				}
			}
		}
		if alreadyTagged {
			continue
		}

		// Apply rules in order
		tagged := false
		rules := tag.Rules
		if tag.Exhaustive {
			rules = append(rules, AskWithPrompt(tag))
		}
		for _, r := range rules {
			selected, err := r(p)
			if err != nil {
				return err
			}
			if selected != "" {
				fmt.Printf("selected option %s for %s\n", selected, Name(p))
				// add membership and add to batch
				p.Memberships = append(p.Memberships, &peoplev1.Membership{
					ContactGroupMembership: &peoplev1.ContactGroupMembership{
						ContactGroupResourceName: groups[selected].ResourceName, // yolo
					},
				})
				batch = append(batch, p)
				tagged = true
				break // stop applying rules
			}
		}
		if !tagged && tag.Exhaustive {
			fmt.Printf("no decision made for %s, skipping...\n", Name(p))
		}

		if len(batch) >= batchSize {
			fmt.Println("saving batch...")
			err = updateAll(srv, batch, "memberships")
			if err != nil {
				return err
			}
			batch = nil
			time.Sleep(2 * time.Second)
		}
	}

	if len(batch) > 0 {
		fmt.Println("saving last batch...")
		err = updateAll(srv, batch, "memberships")
		if err != nil {
			return err
		}
		batch = nil
	}

	return nil
}

func Name(p *peoplev1.Person) string {
	if len(p.Names) > 0 {
		return p.Names[0].DisplayName
	} else {
		return fmt.Sprintf("No Name -- %s", Link(p))
	}
}

func Updated(p *peoplev1.Person) string {
	for _, s := range p.Metadata.Sources {
		if s.Type == "CONTACT" && s.UpdateTime != "" {
			return "updated on " + s.UpdateTime[:10]
		}
	}
	return ""
}

func Note(p *peoplev1.Person) string {
	if len(p.Biographies) > 0 {
		return strings.Split(p.Biographies[0].Value, "\n")[0]
	}
	return ""
}

func Facebook(p *peoplev1.Person) string {
	for _, u := range p.Urls {
		if u.Type == "Facebook" {
			return u.Value
		}
	}
	return ""
}
