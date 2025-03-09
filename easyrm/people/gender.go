package people

import (
	"fmt"
	"sort"

	"github.com/AlecAivazis/survey/v2"
	peoplev1 "google.golang.org/api/people/v1"
)

// AssignGender ensures that every person has a gender. It will ask the user to select
// from "male", "female", or "unspecified".
// NOTE on setting gender with API: values of "male", "female" will map to "male", "female".
// Anything else will map to "other". Once a gender is set, it cannot be removed (but it
// can be changed).
func AssignGender(srv *peoplev1.Service, people []*peoplev1.Person) {

	sort.Slice(people, func(i, j int) bool {
		if len(people[i].Names) == 0 {
			return true
		} else if len(people[j].Names) == 0 {
			return true
		} else {
			return people[i].Names[0].DisplayName < people[j].Names[0].DisplayName
		}
	})

	seen := map[string]string{} // map from name to last-seen gender
	for _, p := range people {
		if len(p.Genders) == 1 && len(p.Names) > 0 {
			seen[p.Names[0].GivenName] = p.Genders[0].Value
		}
	}

	batchSize := 10
	batch := []*peoplev1.Person{}
	save := true
	for _, p := range people {
		if len(p.Genders) >= 1 {
			continue // we're done, no need to ask
		}
		if len(p.Names) < 1 && len(p.Organizations) > 0 {
			fmt.Printf("setting %s's (company) gender to 'other'\n", p.Organizations[0].Name)
			p.Genders = []*peoplev1.Gender{{Value: GenderUnspecified}}
			batch = append(batch, p)
			continue
		}

		ans := GenderUnspecified
		guess := GenderMale
		if g, ok := seen[p.Names[0].GivenName]; ok {
			guess = g
		}
		err := survey.AskOne(&survey.Select{
			Message: fmt.Sprintf("Is %s male/female/skip?", p.Names[0].DisplayName),
			Options: []string{GenderMale, GenderFemale, GenderUnspecified},
			Default: guess,
		}, &ans)
		if err != nil {
			fmt.Println("ERR: " + err.Error())
			save = false
			break
		}
		fmt.Printf("guess is %s. ans is %s\n", seen[p.Names[0].GivenName], ans)

		p.Genders = []*peoplev1.Gender{{Value: ans}}
		batch = append(batch, p)

		if len(batch) >= batchSize {
			fmt.Println("saving batch...")
			err = updateAll(srv, batch, "genders")
			if err != nil {
				fmt.Println("ERR: " + err.Error())
			}
			batch = nil
		}
	}

	if save && len(batch) > 0 {
		fmt.Println("saving last batch...")
		err := updateAll(srv, batch, "genders")
		if err != nil {
			fmt.Println("ERR: " + err.Error())
		}
		batch = nil
	}
}

var GenderMale = "male"
var GenderFemale = "female"
var GenderUnspecified = "unspecified"

func SetGender(srv *peoplev1.Service, person *peoplev1.Person, gender string) error {
	person.Genders = []*peoplev1.Gender{{Value: gender}}
	return update(srv, person, "genders")
}

func ClearGender(srv *peoplev1.Service, person *peoplev1.Person) error {
	person.Genders = []*peoplev1.Gender{}
	return update(srv, person, "genders")
}
