package cleanup

import (
	"easyrm/people"
	"fmt"

	peoplev1 "google.golang.org/api/people/v1"
)

var tags = []*people.Tag{
	{
		Name:        "Attractiveness",
		Description: "How attractive I find the person.",
		Options: []people.Option{
			{"a:no", "not very attractive or cute but too old or untenable"},
			{"a:3", "would bang if easy or because horny; one night stand; not conventionally attractive"},
			{"a:2", "would bang enthusiastically; fuck buddy; nice body, but face is not great or not my type"},
			{"a:1", "pretty face and good enough body; long-term option"},
			{"a:0", "holy shit she's hot"},
			{"a:idk", "i don't know"},
		},
		Exhaustive:          true,
		CreateIfNonExistent: true,
		Applicable: func(p *peoplev1.Person) bool { // Only females
			if len(p.Genders) > 0 && p.Genders[0].Value == people.GenderFemale {
				return true
			}
			return false
		},
	},
	{
		Name:        "Relationship",
		Description: "Whether I should pursue a relationship with them or not.",
		Options: []people.Option{
			{"r:nc-nr", "no chance, no relationship (we don't really know each other and it's hard to get to know each other)"},
			{"r:nc-m", "no chance, she's married or has kids with someone"},
			{"r:nc-ni", "no chance, she's not interested (i tried at least a bit)"},
			{"r:yc-pn", "there's a chance, but it's slim"},
			{"r:yc-pm", "there's a chance.. it's a mayybe (she could be interested)"},
			{"r:yc-py", "there's a chance, and it's likely!"},
			{"r:k", "we kissed"},
			{"r:d", "already had a relationships with them, with intimacy"},
		},
		Exhaustive:          true,
		CreateIfNonExistent: true,
		Applicable: func(p *peoplev1.Person) bool { // Contacts with a:0 or a:1 or a:2
			a0 := "contactGroups/24968e72094023e2"
			a1 := "contactGroups/95a4b638c600e97"
			a2 := "contactGroups/5d42b4388df947a6"
			for _, m := range p.Memberships {
				if m.ContactGroupMembership.ContactGroupResourceName == a0 ||
					m.ContactGroupMembership.ContactGroupResourceName == a1 ||
					m.ContactGroupMembership.ContactGroupResourceName == a2 {
					return true
				}
			}
			return false
		},
	},
	{
		Name:        "Warm",
		Description: "Whether I should stay in touch with this prospect.",
		Options: []people.Option{
			{"Keep warm", "girls I should stay in touch with"},
			{"Keep warmish", "girls I should stay in touch with.. broaded"},
		},
		Exhaustive:          false,
		CreateIfNonExistent: true,
		Applicable: func(p *peoplev1.Person) bool { // Contacts with a:0 or a:1 or a:2
			a0 := "contactGroups/24968e72094023e2"
			a1 := "contactGroups/95a4b638c600e97"
			a2 := "contactGroups/5d42b4388df947a6"
			for _, m := range p.Memberships {
				if m.ContactGroupMembership.ContactGroupResourceName == a0 ||
					m.ContactGroupMembership.ContactGroupResourceName == a1 ||
					m.ContactGroupMembership.ContactGroupResourceName == a2 {
					return true
				}
			}
			return false
		},
		Rules: []people.Rule{func(p *peoplev1.Person) (string, error) {
			a0 := "contactGroups/24968e72094023e2"
			a1 := "contactGroups/95a4b638c600e97"
			a2 := "contactGroups/5d42b4388df947a6"

			unlikely := "contactGroups/2f4ee4900cf3c005"
			maybe := "contactGroups/5bc6f8298e1542a7"
			likely := "contactGroups/2acac30f89104a1d"

			hotness := -1
			likelihood := -1
			for _, m := range p.Memberships {
				if m.ContactGroupMembership.ContactGroupResourceName == a0 {
					hotness = 0
				} else if m.ContactGroupMembership.ContactGroupResourceName == a1 {
					hotness = 1
				} else if m.ContactGroupMembership.ContactGroupResourceName == a2 {
					hotness = 2
				}
				if m.ContactGroupMembership.ContactGroupResourceName == unlikely {
					likelihood = 0
				} else if m.ContactGroupMembership.ContactGroupResourceName == maybe {
					likelihood = 1
				} else if m.ContactGroupMembership.ContactGroupResourceName == likely {
					likelihood = 2
				}
			}

			if hotness == 0 && likelihood != -1 || // any a0 with any chance
				hotness == 1 && likelihood >= 1 || // a1's with some chance (maybe or high)
				hotness == 2 && likelihood == 2 { // a2's with high chance (low effort)
				return "Keep warm", nil
			} else if likelihood != -1 { // any a012 with any chance
				return "Keep warmish", nil
			}
			return "", nil
		}},
	},
	{
		Name:        "Wedding",
		Description: "Whether to invite to wedding or not.",
		Options: []people.Option{
			{"w:no", "do not invite to wedding"},
			{"w:maybe", "maybe invite to wedding"},
			{"w:yes", "invite to wedding"},
			{"w:idk", "i dont know, ask me later again"},
		},
		Exhaustive:          true,
		CreateIfNonExistent: true,
	},
}

func ApplyTags(srv *peoplev1.Service, peeps []*peoplev1.Person) {
	for _, tag := range tags {
		fmt.Printf("Applying tag %s...\n", tag.Name)
		err := people.ApplyTag(srv, peeps, tag)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
