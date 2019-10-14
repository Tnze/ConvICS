package ConvICS

import "time"

// Timetable 是一天的时间表，规定了每天每节课的上课下课时间
type Timetable []struct {
	Start, End time.Time
}

// Schedule 是整个学期的课程安排
type Schedule struct {
}

// Subject 是一门课程
type Subject struct {
}

// Course 是一节课
type Course struct {
}
