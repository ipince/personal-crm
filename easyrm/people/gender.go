package people

import (
	"fmt"
	"sort"

	"github.com/AlecAivazis/survey/v2"
	peoplev1 "google.golang.org/api/people/v1"
)

// AssignGender ensures that every person has a gender. It will ask the user to select
// from "male", "female", or "unspecified".
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

	batchSize := 20
	batch := []*peoplev1.Person{}
	save := true
	for _, p := range people {
		if len(p.Genders) >= 1 {
			continue
		}
		if len(p.Names) < 1 {
			fmt.Printf("skippping %s with no name\n", Link(p))
			continue
		}

		ans := GenderUnspecified
		err := survey.AskOne(&survey.Select{
			Message: fmt.Sprintf("Is %s male/female/skip?", p.Names[0].DisplayName),
			Options: []string{GenderMale, GenderFemale, GenderUnspecified},
		}, &ans)
		if err != nil {
			fmt.Println("ERR: " + err.Error())
			save = false
			break
		}

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
