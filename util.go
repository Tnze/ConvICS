package ConvICS

import (
	"time"
)

func timeMustParse(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

func (t *Timetable) Append(desc string, start, end time.Duration) *Timetable {
	*t = append(*t, struct {
		Description string
		Start, End  time.Duration
	}{
		desc,
		start, end,
	})
	return t
}

func (s Schedule) GetTime(t Timetable, offsetDay, i, d int) (start, end time.Time) {
	ss := s.SemesterStart
	start = ss.AddDate(0, 0, offsetDay).Add(t[i-1].Start)
	end = ss.AddDate(0, 0, offsetDay).Add(t[i-1+d-1].End)
	return
}
