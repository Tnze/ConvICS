package ConvICS

import "time"

func timeMustParse(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

func (t *Timetable) Append(desc string, start, end time.Time) *Timetable {

	*t = append(*t, struct {
		Description string
		Start, End  time.Time
	}{
		desc,
		start, end,
	})
	return t
}
