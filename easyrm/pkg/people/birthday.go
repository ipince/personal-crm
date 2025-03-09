package people

import (
	"fmt"

	peoplev1 "google.golang.org/api/people/v1"
)

type Bday struct {
	Day   int64
	Month int64
	Year  int64 // 0 if no year
}

func (b *Bday) sameDayMonth(o *Bday) bool {
	return b.Day == o.Day && b.Month == o.Month
}

func (b *Bday) hasYear() bool {
	return b.Year != 0
}
func (b *Bday) String() string {
	return fmt.Sprintf("%d/%d/%d", b.Month, b.Day, b.Year)
}

func MaybeSetBirthday(p *peoplev1.Person, bday *Bday) bool {
	if len(p.Birthdays) > 0 && p.Birthdays[0].Date != nil { // should at most =1
		existing := &Bday{
			Day:   p.Birthdays[0].Date.Day,
			Month: p.Birthdays[0].Date.Month,
			Year:  p.Birthdays[0].Date.Year,
		}
		if *bday == *existing || existing.sameDayMonth(bday) && existing.hasYear() && !bday.hasYear() {
			//fmt.Printf("skipping adding bday to %s because it adds no new info\n", Name(p))
			return false
		}
	}

	// If we set as text, it remains as text, even if the data is parseable.
	// Thus, we always set as a Date instead.
	p.Birthdays = []*peoplev1.Birthday{{
		Date: &peoplev1.Date{
			Day:   bday.Day,
			Month: bday.Month,
			Year:  bday.Year,
		}},
	}

	fmt.Printf("set bday %s to %s\n", bday, Name(p))
	return true
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
