package cleanup

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"easyrm/pkg/people"

	peoplev1 "google.golang.org/api/people/v1"
)

func addBirthdays(srv *peoplev1.Service, all []*peoplev1.Person) error {
	birthdaysCsv := "data/facebook-dates-of-birth.csv" // name, year, month, day, facebookID (TODO: some exceptions)
	f, err := os.Open(birthdaysCsv)
	if err != nil {
		return err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	bdays, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	birthdaysByName := map[string]*people.Bday{} // name -> bday
	for _, record := range bdays {
		day, err := strconv.ParseInt(record[3], 10, 0)
		if err != nil {
			fmt.Printf("WARN: failed to parse %s as int (full row is %s). skipping\n", record[3], record)
			continue
		}
		month, err := strconv.ParseInt(record[2], 10, 0)
		if err != nil {
			fmt.Printf("WARN: failed to parse %s as int (full row is %s). skipping\n", record[2], record)
			continue
		}
		year := int64(0)
		if record[1] != "null" {
			year, err = strconv.ParseInt(record[1], 10, 0)
			if err != nil {
				fmt.Printf("WARN: failed to parse %s as int (full row is %s). skipping\n", record[1], record)
				continue
			}
		}
		birthdaysByName[record[0]] = &people.Bday{day, month, year}
	}
	fmt.Printf("read %d birthdays\n", len(birthdaysByName))

	fbUrlsCsv := "data/fb_export_2022-10-14.csv" // url, name
	f, err = os.Open(fbUrlsCsv)
	if err != nil {
		return err
	}
	defer f.Close()

	csvReader = csv.NewReader(f)
	fbUrls, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	fmt.Printf("read %d fb urls\n", len(fbUrls))
	birthdaysByUrl := map[string]*people.Bday{} // url -> bday
	for _, record := range fbUrls {
		name := strings.TrimSpace(record[1])
		if bday, ok := birthdaysByName[name]; ok {
			birthdaysByUrl[record[0]] = bday
		}
	}

	count := 0
	save := []*peoplev1.Person{}
	for _, p := range all {
		added := false
		for _, url := range people.FacebookURLs(p) {
			if bday, ok := birthdaysByUrl[url]; ok {
				count++
				added = people.MaybeSetBirthday(p, bday)
				break // need to break; if a contact has many matching urls (with the same bday), then subsequent calls will make added=false
			}
		}
		if added {
			save = append(save, p)
		}
	}

	fmt.Printf("will add %d birthdays to contacts\n", len(save))

	return people.UpdateAll(srv, save, "birthdays")
}

func birthdays(srv *peoplev1.Service, peeps []*peoplev1.Person) {

	bytes, err := os.ReadFile("data/fb_bday_order.log")
	if err != nil {
		fmt.Println(err.Error())
	}
	links := strings.Split(string(bytes), "\n")

	set := map[string]bool{}
	m := map[string]*peoplev1.Person{}
	for _, l := range links {
		set[l] = true
	}

	for _, p := range peeps {
		for _, u := range people.FacebookURLs(p) {
			if _, found := set[u]; found {
				m[u] = p
			}
		}
	}
	//fmt.Printf("%+v\n", m)

	for _, l := range links {
		if len(m[l].Names) > 0 {
			bday := "???"
			if len(m[l].Birthdays) > 0 {
				bday = m[l].Birthdays[0].Text
			}
			fmt.Printf("%s -> %s\n", m[l].Names[0].DisplayName, bday)
		}
	}
}
