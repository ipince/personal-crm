package cleanup

import (
	"easyrm/people"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"golang.org/x/exp/slices"
	peoplev1 "google.golang.org/api/people/v1"
)

func listPeopleIn(all []*peoplev1.Person, cities []string) {
	for _, p := range all {
		for _, ud := range p.UserDefined {
			if ud.Key == "Current City" && slices.Contains(cities, ud.Value) {
				fmt.Println(people.Name(p))
			}
		}
	}
}

func setCurrentCity(srv *peoplev1.Service, all []*peoplev1.Person, city string) error {
	// select people
	names := []string{}
	byName := map[string]*peoplev1.Person{}
	for _, p := range all {
		if len(p.Names) > 0 {
			name := p.Names[0].DisplayName
			names = append(names, name)
			byName[name] = p
		}
	}

	selection := []string{}
	for {
		selected := ""
		err := survey.AskOne(&survey.Select{
			Message: fmt.Sprintf("Which friends live in %s? (type \"exit\" to exit)", city),
			Options: append(names, "exit"),
		}, &selected)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			fmt.Println("user cancelled out of the prompt. please use \"exit\" instead")
		}
		if selected == "exit" {
			break
		}
		selection = append(selection, selected)
	}

	fmt.Printf("Setting city for %d friends:\n", len(selection))
	key := "Current City"
	save := []*peoplev1.Person{}
	for _, s := range selection {
		p := byName[s]
		fmt.Println(people.Name(p))

		found := false
		for _, ud := range p.UserDefined {
			if ud.Key == key { // Replace, if different
				found = true
				if ud.Value != city {
					fmt.Printf("%s: %s -> %s\n", people.Name(p), ud.Value, city)
					ud.Value = city
					save = append(save, p)
				}
				continue
			}
		}

		if !found { // Add if not found
			p.UserDefined = append(p.UserDefined, &peoplev1.UserDefined{
				Key:   key,
				Value: city,
			})
			save = append(save, p)
		}
	}

	return people.UpdateAll(srv, save, "userDefined")
}
