package ConvICS

import (
	"time"

	"github.com/google/uuid"
)

// Timetable 是一天的时间表，规定了每天每节课的上课下课时间
type Timetable []struct {
	Description string
	Start, End  time.Duration
}

// Schedule 是整个学期的课程安排
type Schedule struct {
	SemesterStart time.Time // 学期开始的第一天，用于计算周数
	TotalWeeks    int       // 总周数
	Subjects      map[uuid.UUID][]Subject
}

// Subject 是一门课程
type Subject struct {
	Name         string
	Teacher      string
	Location     string
	Weekday      time.Weekday
	Start, End   int // 第[Start, End]周，从1开始
	CStart, CEnd int // 第[CStart, CEnd]节课，从1开始
}
