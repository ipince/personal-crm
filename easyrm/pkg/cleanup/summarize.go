package cleanup

import (
	"easyrm/contact"
	"easyrm/people"
	"fmt"
	"sort"
	"strings"

	"golang.org/x/exp/slices"
	peoplev1 "google.golang.org/api/people/v1"
)

func Summarize(srv *peoplev1.Service, all []*peoplev1.Person) error {
	summarizeGenders(all)
	summarizeBirthdays(all)
	summarizeAddresses(all)
	summarizeCustomKeyPairs(all)
	summarizeNotes(all)
	summarizeUrls(all)
	return summarizeTags(srv, all)
}

func summarizeTags(srv *peoplev1.Service, all []*peoplev1.Person) error {
	groups, err := srv.ContactGroups.List().Do()
	if err != nil {
		return err
	}

	namesByID := map[string]string{} // gid -> group name (aka option name)
	for _, g := range groups.ContactGroups {
		gid := g.ResourceName
		namesByID[gid] = g.Name
	}

	optionNames := []string{}
	totals := map[string]map[string]int{} // Tag -> Option -> Count
	for _, t := range tags {
		for _, opt := range t.Options {
			optionNames = append(optionNames, opt.Name)
		}
		totals[t.Name] = map[string]int{}
		eligible := 0
		for _, p := range all {
			if t.AppliesTo(p) {
				eligible++
				for _, m := range p.Memberships {
					gid := m.ContactGroupMembership.ContactGroupResourceName
					gName := namesByID[gid]
					for _, opt := range t.Options {
						if gName == opt.Name {
							totals[t.Name][opt.Name] += 1
							break
						}
					}
				}
			}
		}
		totals[t.Name]["eligible"] = eligible
	}

	totals["Other"] = map[string]int{}
	totals["Other"]["eligible"] = len(all)
	for _, p := range all {
		for _, m := range p.Memberships {
			gid := m.ContactGroupMembership.ContactGroupResourceName
			gName := namesByID[gid]
			if gName == "" {
				continue // myContacts
			} else if slices.Contains(optionNames, gName) {
				continue // already counted elsewhere
			}
			totals["Other"][gName] += 1
		}
	}

	fmt.Println("Tags:")
	for t, counts := range totals {
		denom := counts["eligible"]
		fmt.Printf("  %s (%d)\n", t, denom)

		names := []string{}
		for name, _ := range counts {
			names = append(names, name)
		}
		sort.Slice(names, func(i, j int) bool { return names[i] < names[j] })

		for _, name := range names {
			if name == "eligible" {
				continue
			}
			fmt.Printf("  - %s: %d (%.2f%%)\n", name, counts[name], float64(counts[name])*100/float64(denom))
		}
	}

	return nil
}

func summarizeUrls(all []*peoplev1.Person) {
	counts := map[string]int{}
	for _, p := range all {
		urls := 0
		for _, u := range p.Urls {
			if u.Metadata.Source.Type != "CONTACT" {
				continue
			}

			urls++
			counts[u.Type] += 1
		}

		if urls == 0 {
			counts["(none)"] += 1
		}
	}

	fmt.Println("URLs:")
	for k, v := range counts {
		fmt.Printf("  - %s: %d (%.2f%%)\n", k, v, float64(v)*100/float64(len(all)))
	}
}

func summarizeBirthdays(all []*peoplev1.Person) {
	counts := map[string]int{}
	for _, p := range all {
		// Assume len(birthdays) == 0
		if len(p.Birthdays) == 0 {
			counts["(none)"] += 1
			continue
		}
		b := p.Birthdays[0]
		if b.Metadata.Source.Type != "CONTACT" {
			fmt.Println(contact.New(p).ShortString())
			continue
		}

		if b.Text != "" {
			counts["Text"] += 1
		} else if b.Date != nil {
			if b.Date.Year == 0 {
				counts["Day-Only"] += 1
			} else {
				counts["Full"] += 1
			}
		}
	}

	fmt.Println("Birthdays:")
	for k, v := range counts {
		fmt.Printf("  - %s: %d (%.2f%%)\n", k, v, float64(v)*100/float64(len(all)))
	}
}

func summarizeAddresses(all []*peoplev1.Person) {
	counts := map[string]int{}
	for _, p := range all {
		addrs := 0
		for _, a := range p.Addresses {
			if a.Metadata.Source.Type != "CONTACT" {
				fmt.Println("contact...")
				continue
			}

			addrs++
			counts[a.Type] += 1
		}

		if addrs == 0 {
			counts["(none)"] += 1
		}
	}

	fmt.Println("Addresses:")
	for k, v := range counts {
		fmt.Printf("  - %s: %d (%.2f%%)\n", k, v, float64(v)*100/float64(len(all)))
	}
}

func summarizeCustomKeyPairs(all []*peoplev1.Person) {
	counts := map[string]int{}
	for _, p := range all {
		for _, ud := range p.UserDefined {
			counts[ud.Key] += 1
			if strings.Contains(ud.Value, "Mexico") {
				fmt.Printf("%s -> %s\n", ud.Value, people.Name(p))
			}
		}

		if len(p.UserDefined) == 0 {
			counts["(none)"] += 1
		}
	}

	fmt.Println("Key Pairs:")
	for k, v := range counts {
		fmt.Printf("  - %s: %d (%.2f%%)\n", k, v, float64(v)*100/float64(len(all)))
	}
}

func summarizeNotes(all []*peoplev1.Person) {
	counts := map[string]int{}
	for _, p := range all {
		for _, n := range p.Biographies {
			if n.Value != "" {
				counts["notes"] += 1
			}
		}

		if len(p.Biographies) == 0 {
			counts["(none)"] += 1
		}
	}

	fmt.Println("Notes:")
	for k, v := range counts {
		fmt.Printf("  - %s: %d (%.2f%%)\n", k, v, float64(v)*100/float64(len(all)))
	}
}

func summarizeGenders(all []*peoplev1.Person) {
	hist := map[string]int{}
	for _, p := range all {
		val := "none"
		if len(p.Genders) == 1 {
			val = p.Genders[0].Value
		} else if len(p.Genders) > 1 {
			val = "multi"
		}

		if val == "none" {
			fmt.Printf("this person has no gender:\n%s\n", contact.New(p).ShortString())
		}

		if _, ok := hist[val]; !ok {
			hist[val] = 1
		} else {
			hist[val] += 1
		}
	}
	fmt.Println("Genders:")
	for k, v := range hist {
		fmt.Printf("  - %s: %d (%.2f%%)\n", k, v, float64(v)*100/float64(len(all)))
	}
}
