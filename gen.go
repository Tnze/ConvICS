package ConvICS

import (
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/google/uuid"
)

func (s Schedule) ToICS(t Timetable) []byte {
	cal := ics.NewCalendar()
	cal.SetProductId("github.com/Tnze/ConvICS")

	for _, subject := range s.Subjects {
		for _, course := range subject {
			for i := course.Start; i <= course.End; i++ {
				event := cal.AddEvent(uuid.New().String())
				event.SetCreatedTime(time.Now())
				event.SetDtStampTime(time.Now())
				event.SetModifiedAt(time.Now())
				s, e := s.GetTime(t, (i-1)*7+int(course.Weekday), course.CStart, course.CEnd)
				event.SetStartAt(s)
				event.SetEndAt(e)
				event.SetSummary(course.Name)
				event.SetLocation(course.Location)
				event.SetDescription(course.Teacher)
				//event.SetURL("https://URL/")
				//event.SetOrganizer("sender@domain", ics.WithCN("This Machine"))
				//event.AddAttendee("reciever or participant", ics.CalendarUserTypeIndividual, ics.ParticipationStatusNeedsAction, ics.ParticipationRoleReqParticipant, ics.WithRSVP(true))
			}
		}
	}
	return []byte(cal.Serialize())
}
