package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"time"
)

func main() {
	fileName := os.Args[1] //文件名
	if path.Ext(fileName) != ".yml" {
		panic(fmt.Errorf("应该是一个yml文件！"))
	}

	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(fmt.Errorf("读文件%s失败: %v", fileName, err))
	}
	var s Schedule
	yaml.Unmarshal(file, &s)
	//fmt.Print(s)

	var ics string
	ics += icsHand()
	for i := 0; i < len(s.Classes); i++ {
		st, err := time.Parse("20060102T150405-07", s.FirstWeek+"T"+s.Timetable[s.Classes[i].Time][0])
		if err != nil {
			panic(err)
		}
		en, err := time.Parse("20060102T150405-07", s.FirstWeek+"T"+s.Timetable[s.Classes[i].Time][1])
		if err != nil {
			panic(err)
		}
		for j := s.Classes[i].Week[0] - 1; j < s.Classes[i].Week[1]; j++ {
			start := st.AddDate(0, 0, 7*j+s.Classes[i].Day)
			end := en.AddDate(0, 0, 7*j+s.Classes[i].Day)
			ics += icsEvent(start, end, s.Classes[i].Name)
		}
	}
	ics += icsEnd()

	icsFile, err := os.OpenFile(fileName+".ics", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	defer icsFile.Close()

	if err != nil {
		panic(fmt.Errorf("创建文件%s失败: %v", fileName+".ics", err))
	}

	if _, err := icsFile.WriteString(ics); err != nil {
		panic(fmt.Errorf("写入文件失败: %v", err))
	}
}

type Schedule struct {
	Classes   []Class              `yaml:"classes"`
	Timetable map[string][2]string `yaml:"schedule"`
	FirstWeek string               `yaml:"firstWeek"`
}

type Class struct {
	Time     string `yaml:"time"`
	Day      int    `yaml:"day"`
	Name     string `yaml:"name"`
	Teacher  string `yaml:"teacher"`
	Week     [2]int `yaml:"week"`
	Location string `yaml:"location"`
}

func icsHand() string {
	return `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Tnze//YAML-iCalendar v1.0//CN
`
}

func icsEnd() string {
	return "END:VCALENDAR"
}

func icsEvent(start, end time.Time, summary string) string {
	return fmt.Sprintf("BEGIN:VEVENT\nUID:%s\nDTSTAMP:19970714T170000Z\nDTSTART:%s\nDTEND:%s\nSUMMARY:%s\nEND:VEVENT\n",
		uuid.Must(uuid.NewV4()),
		start.UTC().Format("20060102T150405Z"),
		end.UTC().Format("20060102T150405Z"),
		summary,
	)
}
