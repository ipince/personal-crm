package people

import (
	"fmt"
	"regexp"
	"strings"

	peoplev1 "google.golang.org/api/people/v1"
)

func Normalize(srv *peoplev1.Service, person *peoplev1.Person) error {
	return normalizePhoneNumbers(srv, person)
}

// NormalizeUrls normalizes the URLs of a person
func NormalizeUrls(srv *peoplev1.Service, person *peoplev1.Person) error {
	domains := map[string]string{
		"facebook.com":  "Facebook",
		"linkedin.com":  "LinkedIn",
		"instagram.com": "Instagram",
		"youtube.com":   "YouTube",
	}
	remove := "plus.google.com"

	rmIndex := -1
	dirty := false

	for i, u := range person.Urls {
		if u.Metadata.Source.Type != "CONTACT" {
			continue
		}

		if strings.Contains(u.Value, remove) {
			rmIndex = i
			dirty = true
			continue
		}

		known := false
		for domain, label := range domains {
			if strings.Contains(u.Value, domain) {
				known = true
				if u.Type != label {
					fmt.Printf("%s: %s (%s -> %s)\n", Name(person), u.Value, u.Type, label)
					u.Type = label
					dirty = true
					break
				}
			}
		}
		if !known && u.Type != "other" {
			fmt.Printf("%s: %s (%s -> %s)\n", Name(person), u.Value, u.Type, "other")
			u.Type = "other"
			dirty = true
		}
	}

	if rmIndex != -1 { // remove...
		person.Urls = append(person.Urls[:rmIndex], person.Urls[rmIndex+1:]...)
		dirty = true
	}

	if dirty {
		return update(srv, person, "urls")
	}

	return nil // nothing to update
}

func normalizePhoneNumbers(srv *peoplev1.Service, person *peoplev1.Person) error {
	msgs := []string{}
	for i, n := range person.PhoneNumbers {
		normalized := transformUSNumber(n.Value)
		if n.Value != normalized {
			msgs = append(msgs, fmt.Sprintf("transforming %s -> %s\n", n.Value, normalized))
			person.PhoneNumbers[i].Value = normalized
		}
	}
	if len(msgs) > 0 {
		fmt.Println(Link(person))
		fmt.Println(strings.Join(msgs, "\n"))

		return update(srv, person, "phoneNumbers")
	}

	return nil // nothing to update
}

func transformUSNumber(number string) string {
	// Omit numbers that start with +, but are not +1. No other countries
	// start with +1 (well, except for territories with codes like +1-###).
	// https://countrycode.org/
	if strings.HasPrefix(number, "+") {
		if len(number) > 1 {
			if string(number[1]) != "1" {
				// skip; it already has a country code for another country.
				return number
			}
		} else { // malformed number
			fmt.Println("Ignoring malformed number: " + number)
			return number
		}
	}

	// Remove everything that isn't a number.
	cleared := clearString(number)
	if len(cleared) == 11 && strings.HasPrefix(cleared, "1") {
		// Divide in chunks of 3-3-4.
		return formatAsUSNumber(cleared[1:])
	} else if len(cleared) == 10 {
		// Probably non-US/Canada
		return formatAsUSNumber(cleared)
	} else { // must be non-US/Canada
		return number
	}
}

var nonNumericRegex = regexp.MustCompile(`[^0-9]+`)

func clearString(str string) string {
	return nonNumericRegex.ReplaceAllString(str, "")
}

func formatAsUSNumber(number string) string {
	// number has 10 digits
	return fmt.Sprintf("+1 %s-%s-%s", number[0:3], number[3:6], number[6:])
}
