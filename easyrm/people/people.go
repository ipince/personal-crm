package people

import (
	"errors"
	"fmt"
	"log"
	"strings"

	peoplev1 "google.golang.org/api/people/v1"
)

const TestPersonID = "people/c4490748168429910297"

var PeopleFields = []string{

	// Basic
	"names",
	"nicknames",
	"birthdays", // set for all.
	"photos",

	// Contact
	"addresses", // Physical addresses
	"emailAddresses",
	"phoneNumbers",
	"userDefined", // Custom: Current City, Past City

	// Work -- settable both in Contacts and Profile?
	"organizations",
	"userDefined", // Custom: Current Company, Past Company

	// Important dates / events
	"events", // maybe day-we-met

	// Links
	"imClients",
	"urls", // facebook url (or others, can have a label)

	// Relationships? -- I probably won't really use
	"relations",

	// Groups
	"memberships", // labels...at least MyContacts.
	// Closeness: Close, Friend, Acquientance
	// Circle:
	// ... what if I just consider them tags? and use custom fields? we'll see.

	// Notes
	"biographies", // aka "notes"

	// Metadata -- last updated? created?
	"metadata",

	// Unused so far. Not settable in Google Contacts UI... where does it come from?
	"ageRanges",
	"calendarUrls",
	"clientData",
	// "coverPhotos", // REMOVED: mostly default photo from Google+ (source=PROFILE)
	"externalIds",
	// "genders",  // TODO: Maybe use? TEMP
	"interests",
	"locales",
	"locations", // TODO decide what to do
	"miscKeywords",
	"occupations", // TODO: decide what to do
	"sipAddresses",
	"skills",
}

func Get(srv *peoplev1.Service, id string) (*peoplev1.Person, error) {
	return srv.People.Get(id).
		PersonFields(allFields()).Do()
}

func SetBirthdate(srv *peoplev1.Service, id string, day, month int64, year *int64) error {

	person, err := Get(srv, id)
	if err != nil {
		return err // TODO
	}
	printBirthdays(person)

	// If we set as text, it remains as text, even if the data is parseable.
	// Thus, we always set as a Date instead.
	dob := &peoplev1.Birthday{
		Date: &peoplev1.Date{
			Day:   day,
			Month: month,
		},
	}
	if year != nil {
		dob.Date.Year = *year
	}
	person.Birthdays[0] = dob

	_, err = srv.People.UpdateContact(id, person).
		UpdatePersonFields("birthdays").Do()
	return err
}

func ListAll(srv *peoplev1.Service) ([]*peoplev1.Person, error) {
	const pageSize = 1000
	r, err := srv.People.Connections.List("people/me").
		PageSize(pageSize).
		PersonFields(allFields()).Do()
	if err != nil {
		return nil, err
	}

	people := make([]*peoplev1.Person, 0, r.TotalItems)
	for {
		people = append(people, r.Connections...)
		if int64(len(people)) >= r.TotalItems {
			break
		}

		r, err = srv.People.Connections.List("people/me").
			PageSize(pageSize).
			PageToken(r.NextPageToken).
			PersonFields(allFields()).Do()
		if err != nil {
			return nil, err
		}
	}

	return people, nil
}

func List(srv *peoplev1.Service, limit int) {
	r, err := srv.People.Connections.List("people/me").
		PageSize(int64(limit)).
		PersonFields("names,emailAddresses").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve people. %v", err)
	}
	if len(r.Connections) > 0 {
		fmt.Print("List 10 connection names:\n")
		for _, c := range r.Connections {
			names := c.Names
			if len(names) > 0 {
				name := names[0].DisplayName
				fmt.Printf("%s\n", name)
			}
		}
	} else {
		fmt.Print("No connections found.")
	}
}

// Validate prints any non-conforming fields found in the given person
func Validate(person *peoplev1.Person) {
	var errs []error
	err := validateNames(person)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		fmt.Printf("Failed to validate person %s\n", link(person))
		for _, e := range errs {
			fmt.Println(e)
		}
	}
	//printBirthdays(person)
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

func printBirthdays(person *peoplev1.Person) {
	fmt.Printf("Found %d birthdays in person %s: ", len(person.Birthdays), person.Names[0].DisplayName)
	for _, b := range person.Birthdays {
		if b.Date != nil {
			fmt.Printf("date %d/%d/%d, ", b.Date.Day, b.Date.Month, b.Date.Year)
		} else {
			fmt.Printf("text %s, ", b.Text)
		}
	}
	fmt.Println()
}

func link(person *peoplev1.Person) string {
	parts := strings.Split(person.ResourceName, "/")
	return "https://contacts.google.com/person/" + parts[1]
}

func allFields() string {
	return strings.Join(PeopleFields, ",")
}