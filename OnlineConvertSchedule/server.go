package main

import (
	// "encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/satori/go.uuid"
	// "html"
	"io"
	// "io/ioutil"
	"net/http"
	// "os"
	"log"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/convert_schedule", scheduleParser)
	http.HandleFunc("/upload", uploadDocument)
	http.ListenAndServe(":1308", nil)
	// f, err := os.Open("00000000000_XXX(学生课表).xls.exe")
	// if err != nil {
	// 	panic(err)
	// }
	// var data fullTable
	// data.parse(f)
	// j, err := json.Marshal(data)
	// fmt.Println(string(j), err)
	// ioutil.WriteFile("out.ics", []byte(data.toICS()), 0666)
}

func uploadDocument(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	rw.Write([]byte(`<html>

		<head></head>
		
		<body>
			<form action="/convert_schedule" method="POST" enctype="multipart/form-data">
				<input type="file" name="sch" />
				<input type="submit" value="Convert" />
			</form>
		</body>
		
		</html>`))
}

func scheduleParser(rw http.ResponseWriter, r *http.Request) {
	reader := http.MaxBytesReader(rw, r.Body, 100000)
	defer reader.Close()
	data := new(fullTable)
	data.parse(reader)
	if data.Error != nil {
		log.Println(data.Error)
		rw.WriteHeader(505)
		return
	}

	rw.Header().Set("Content-type", "application/octet-stream")
	rw.Header().Set("Content-Disposition", "attachment;filename=yours.ics")
	// log.Println(data.toICS())
	if _, err := rw.Write([]byte(data.toICS())); err != nil {
		log.Println(err)
	}
}

type fullTable struct {
	Info struct {
		Semester string //学期
		Number   string //学号
		Class    string //班级
		Name     string //姓名
		Credit   string //学分
	}
	Curriculums []curri //20教学周，一个星期7天，每天11节课，允许多个课程排在同一时间
	Error       error
}

type curri struct {
	N, M                 int    //第n周-第m周,
	Day, Start, Duration int    //星期几，开始时间，持续时间
	Teacher              string //教师姓名
	Name                 string //课程名称(序号)
	Loc                  string //教室
}

func (f *fullTable) parse(r io.Reader) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		f.Error = err
		return
	}
	doc.Find("table").Each(func(i int, s *goquery.Selection) {
		switch i {
		case 0: //课表信息info
			n := s.Find("tr > td")
			n.Each(func(i int, s *goquery.Selection) {
				switch i {
				case 1:
					f.Info.Semester = s.Text()
				case 2:
					field := strings.Fields(s.Text())

					f.Info.Number = strings.TrimPrefix(field[0], "学号:")
					f.Info.Name = strings.TrimPrefix(field[1], "学生姓名:")
					f.Info.Class = field[3]
					f.Info.Credit = strings.TrimPrefix(field[4], "总学分:")
				}
			})
		case 1: //课表curriculums
			s.Find("tr").EachWithBreak(func(n int, s *goquery.Selection) bool {
				s.Find("td").Each(func(j int, s *goquery.Selection) {
					if n == 0 {
						// fmt.Println(strings.TrimSpace(s.Text()))
					} else if j < 7 && n-1 < 11 { //一星期不超过7天，一天不超过11节课
						// fmt.Println(i-1, j)
						fe := strings.Split(s.Text(), "\n")

						for i := 0; i+5 < len(fe); i += 5 {
							var c curri
							c.Start = n
							c.Day = j
							fmt.Sscan(s.AttrOr("rowspan", "2"), &c.Duration) //课长
							//单个课程占5行，把数据读取到结构体内
							//0
							//1 课程名 (编号)
							//2 (教师)
							//3
							//4 (n-m 地址)
							c.Name = strings.TrimSpace(fe[i+1])
							c.Teacher = strings.Trim(fe[i+2], "( )")
							plan := strings.Trim(fe[i+4], "( )") + ")"
							if _, err := fmt.Sscanf(plan, "%d-%d %s", &c.N, &c.M, &c.Loc); err != nil {
								fmt.Sscanf(plan, "%d %s", &c.N, &c.Loc)
							}

							// fmt.Println(c)
							f.Curriculums = append(f.Curriculums, c)
						}
					}
				})
				return i < 11
			})
		case 2:
			// s.Find("tr").Each(func(i int, s *goquery.Selection) {
			// 	s.Find("td").Each(func(i int, s *goquery.Selection) {
			// 		fmt.Println(s.Text())
			// 	})
			// })
		}
	})

	return
}

var (
	//SchoolDay 定义学期开始第一天
	SchoolDay = time.Date(2019, 3, 3, 0, 0, 0, 0, time.Local)
	//Schedule 十一节课的时间表
	Schedule = [11]struct{ start, end time.Duration }{
		{time.Hour * 8, time.Hour*8 + time.Minute*50},
		{time.Hour * 9, time.Hour*9 + time.Minute*50},
		{time.Hour*10 + time.Minute*10, time.Hour * 11},
		{time.Hour*11 + time.Minute*10, time.Hour * 12},
		{time.Hour * 14, time.Hour*14 + time.Minute*50},
		{time.Hour * 15, time.Hour*15 + time.Minute*50},
		{time.Hour*16 + time.Minute*10, time.Hour * 17},
		{time.Hour*17 + time.Minute*10, time.Hour * 18},
		{time.Hour*19 + time.Minute*30, time.Hour*20 + time.Minute*20},
		{time.Hour*20 + time.Minute*30, time.Hour*21 + time.Minute*20},
		{time.Hour*21 + time.Minute*30, time.Hour*22 + time.Minute*20},
	}
	//DayName 是所有星期的简写
	DayName = [7]string{"SU", "MO", "TU", "WE", "TH", "FR", "SA"}
)

func (f *fullTable) toICS() string {
	sb := new(strings.Builder)
	fmt.Fprintln(sb, "BEGIN:VCALENDAR")
	fmt.Fprintln(sb, "VERSION:2.0")
	fmt.Fprintln(sb, "PRODID:-//Tnze//YAML-iCalendar v1.0//CN")
	for _, v := range f.Curriculums {
		day := SchoolDay.AddDate(0, 0, v.Day+(v.N-1)*7)
		start := day.Add(Schedule[v.Start-1].start)
		end := day.Add(Schedule[v.Start-2+v.Duration].end)
		fmt.Fprintf(sb, "BEGIN:VEVENT\nUID:%s\nDTSTAMP:%s\nDTSTART:%s\nDTEND:%s\nSUMMARY:%s\nLOCATION:%s\nDESCRIPTION:%s\n",
			uuid.Must(uuid.NewV4()),
			time.Now().Format("20060102T150405Z"),
			start.UTC().Format("20060102T150405Z"),
			end.UTC().Format("20060102T150405Z"),
			v.Name,
			v.Loc,
			v.Teacher,
		)
		if v.M != 0 {
			fmt.Fprintf(sb, "RRULE:FREQ=WEEKLY;INTERVAL=1;BYDAY=%s;COUNT=%d\n", DayName[v.Day], 1+v.M-v.N)
		}
		fmt.Fprintln(sb, "END:VEVENT")
	}
	fmt.Fprintln(sb, "END:VCALENDAR")
	return sb.String()
}
