package people

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/davecgh/go-spew/spew"
	peoplev1 "google.golang.org/api/people/v1"
)

// Validate prints any non-conforming fields found in the given person
func Validate(person *peoplev1.Person) {
	if err := validateNames(person); err != nil {
		fmt.Printf("Name validation error for %s: %v\n", Name(person), err)
	}
	if err := validateBirthdays(person); err != nil {
		fmt.Printf("Birthday validation error for %s: %v\n", Name(person), err)
	}
}

// ValidateNames validates the names of a person
func ValidateNames(person *peoplev1.Person) error {
	return validateNames(person)
}

// ValidateBirthdays validates the birthdays of a person
func ValidateBirthdays(person *peoplev1.Person) error {
	return validateBirthdays(person)
}

func validateNames(person *peoplev1.Person) error {
	if len(person.Names) == 0 {
		if len(person.Organizations) != 1 {
			// Only allow no name for a business contact. A business contact
			// only has one organization (the name of the business).
			return errors.New("no name found")
		}
	} else if len(person.Names) > 1 {
		displayNames := []string{}
		for _, n := range person.Names {
			displayNames = append(displayNames, n.DisplayName)
		}
		return errors.New("too many names found (expected 1): " + strings.Join(displayNames, ","))
	}
	return nil
}

func validateBirthdays(person *peoplev1.Person) error {
	if len(person.Birthdays) > 1 {
		log.Printf("%+v", person.Birthdays)
		spew.Dump(person.Birthdays)
		return errors.New("too many birthdays")
	} else if len(person.Birthdays) == 1 {
		if person.Birthdays[0].Date == nil {
			return errors.New("found unstructured birthday")
		}
	}
	return nil
}
