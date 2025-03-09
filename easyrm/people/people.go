package people

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	peoplev1 "google.golang.org/api/people/v1"
)

const TestPersonID = "people/c4490748168429910297"
const NewlyAdded = "people/c5197444844658991492"

var PeopleFields = []string{

	// Basic
	"names",
	"nicknames",
	"birthdays", // set for all.
	"genders",   // TODO: Maybe use? TEMP
	"photos",

	// Contact
	"emailAddresses",
	"phoneNumbers",
	"addresses",   // Physical addresses
	"userDefined", // Custom: Current City, Past City

	// Work -- settable both in Contacts and Profile?
	"organizations",
	"userDefined", // Custom: Current Company, Past Company

	// Important dates / events
	"events", // maybe day-we-met

	// Links
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
	"imClients", // Deprecated by Google Contacts
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

// TODO: remove
var stagingGroupID = "309a2c4c8b9dede5"

func Insert(srv *peoplev1.Service, fullName, fbURL string) error {
	person := &peoplev1.Person{
		Names: []*peoplev1.Name{
			{
				UnstructuredName: fullName,
			},
		},
		Urls: []*peoplev1.Url{
			{
				Type:  "Facebook",
				Value: fbURL,
			},
		},
		Memberships: []*peoplev1.Membership{
			{
				ContactGroupMembership: &peoplev1.ContactGroupMembership{
					ContactGroupId:           stagingGroupID,
					ContactGroupResourceName: fmt.Sprintf("contactGroups/%s", stagingGroupID),
				},
			},
		},
	}

	_, err := srv.People.CreateContact(person).Do()
	if err != nil {
		return err
	}

	return nil
}

func UpdateAll(srv *peoplev1.Service, people []*peoplev1.Person, fields string) error {
	return updateAll(srv, people, fields)
}

var batchSize = 200

func updateAll(srv *peoplev1.Service, people []*peoplev1.Person, fields string) error {
	for i := 0; i < len(people); i += batchSize {
		end := i + batchSize

		if end > len(people) {
			end = len(people)
		}
		err := updateBatch(srv, people[i:end], fields)
		if err != nil {
			return err
		}

		time.Sleep(1 * time.Second)
	}
	return nil
}

func updateBatch(srv *peoplev1.Service, batch []*peoplev1.Person, fields string) error {
	if len(batch) > batchSize {
		return errors.New("batch is too large")
	}
	fmt.Printf("writing batch of size %d...\n", len(batch))

	toUpdate := map[string]peoplev1.Person{}
	for _, p := range batch {
		toUpdate[p.ResourceName] = *p
	}
	_, err := srv.People.BatchUpdateContacts(&peoplev1.BatchUpdateContactsRequest{
		Contacts:   toUpdate,
		UpdateMask: fields,
	}).Do()
	return err
}

func update(srv *peoplev1.Service, person *peoplev1.Person, fields string) error {
	_, err := srv.People.UpdateContact(person.ResourceName, person).
		UpdatePersonFields(fields).Do()
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

func Print(person *peoplev1.Person) {
	j, err := person.MarshalJSON()
	if err != nil {
		fmt.Println(err)
	}
	y, err := yaml.JSONToYAML(j)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(Link(person))
	fmt.Println(string(y))
}

func Link(person *peoplev1.Person) string {
	parts := strings.Split(person.ResourceName, "/")
	return "https://contacts.google.com/person/" + parts[1]
}

func allFields() string {
	return strings.Join(PeopleFields, ",")
}
