package contact

import (
	"fmt"
	"strconv"
	"strings"

	peoplev1 "google.golang.org/api/people/v1"
)

// Contact represents a contact. The type presents a better interface for accessing and
// manipulating a contact that the Google Contacts representation (Person) does.
type Contact struct {
	person *peoplev1.Person
}

func New(p *peoplev1.Person) *Contact {
	return &Contact{
		person: p,
	}
}

// TODO: date added, last updated, labels
var shortFmt = `
names: %s
emails: %s
phones: %s
orgs: %s
urls: %s
genders: %s
birthdays: %s
addresses: %s
cities: %s
notes: %s
`

func (c *Contact) ShortString() string {
	return fmt.Sprintf(
		shortFmt,
		strings.Join(c.names(), ", "),
		strings.Join(c.emails(), ", "),
		strings.Join(c.phones(), ", "),
		strings.Join(c.orgs(), ", "),
		strings.Join(c.urls(), " | "),
		strings.Join(c.genders(), ", "),
		strings.Join(c.birthdays(), ", "),
		strings.Join(c.addresses(), ";"),
		strings.Join(c.cities(), ","),
		strings.Join(c.notes(), ","),
	)
}

func (c *Contact) Name() string {
	for _, v := range c.person.Names {
		return v.DisplayName // display only first name for now
	}
	return "<no-name>"
}

func (c *Contact) FirstAndLastName() string {
	for _, v := range c.person.Names {
		firsts := strings.Split(v.GivenName, " ")
		lasts := strings.Split(v.FamilyName, " ")
		return fmt.Sprintf("%s %s", firsts[0], lasts[0])
	}
	return "<no-name>"
}

func (c *Contact) FirstEmail() string {
	emails := c.emails()
	if len(emails) > 0 {
		return emails[0]
	}
	return ""
}

func (c *Contact) FirstPhone() string {
	phones := c.phones()
	if len(phones) > 0 {
		return phones[0]
	}
	return ""
}

func (c *Contact) Notes() string {
	for _, v := range c.person.Biographies {
		return v.Value // display only first name for now
	}
	return ""
}

func (c *Contact) names() []string {
	vals := []string{}
	for _, v := range c.person.Names {
		vals = append(vals, v.DisplayName)
	}
	return vals
}

func (c *Contact) emails() []string {
	vals := []string{}
	for _, v := range c.person.EmailAddresses {
		vals = append(vals, v.Value)
	}
	return vals
}

func (c *Contact) phones() []string {
	vals := []string{}
	for _, v := range c.person.PhoneNumbers {
		vals = append(vals, v.Value)
	}
	return vals
}
func (c *Contact) orgs() []string {
	vals := []string{}
	for _, v := range c.person.Organizations {
		vals = append(vals, v.Name)
	}
	return vals
}
func (c *Contact) urls() []string {
	vals := []string{c.Link()}
	for _, v := range c.person.Urls {
		vals = append(vals, fmt.Sprintf("%s (%s)", v.Value, v.Type))
	}
	return vals
}
func (c *Contact) Link() string {
	parts := strings.Split(c.person.ResourceName, "/")
	return "https://contacts.google.com/person/" + parts[1]
}
func (c *Contact) genders() []string {
	vals := []string{}
	for _, v := range c.person.Genders {
		if len(c.person.Genders) > 1 {
			vals = append(vals, fmt.Sprintf("%s (%s)", v.Value, v.Metadata.Source.Type))
		} else {
			vals = append(vals, v.Value)
		}
	}
	return vals
}
func (c *Contact) birthdays() []string {
	vals := []string{}
	for _, v := range c.person.Birthdays {
		bday := fmt.Sprintf("text %s", v.Text)
		if v.Date != nil {
			bday = fmt.Sprintf("date %d/%d", v.Date.Day, v.Date.Month)
			if v.Date.Year != 0 {
				bday += "/" + strconv.FormatInt(v.Date.Year, 10)
			}
		}
		vals = append(vals, bday)
	}
	return vals
}
func (c *Contact) addresses() []string {
	vals := []string{}
	for _, v := range c.person.Addresses {
		vals = append(vals, fmt.Sprintf("%s, %s (%s)", v.StreetAddress, v.City, v.Type))
	}
	return vals
}
func (c *Contact) cities() []string {
	vals := []string{}
	for _, v := range c.person.UserDefined {
		if strings.Contains(v.Key, "City") {
			vals = append(vals, fmt.Sprintf("%s (%s)", v.Value, v.Key))
		}
	}
	return vals
}
func (c *Contact) notes() []string {
	vals := []string{}
	for _, v := range c.person.Biographies {
		vals = append(vals, v.Value)
	}
	return vals
}
