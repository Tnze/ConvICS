// 解析教务网上导出的课程表
package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/Tnze/ConvICS"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

var timetable ConvICS.Timetable
var secondsEastOfUTC = int((8 * time.Hour).Seconds())
var beijing = time.FixedZone("Beijing Time", secondsEastOfUTC)
var schedule = ConvICS.Schedule{
	SemesterStart: time.Date(2019, 9, 1, 0, 0, 0, 0, beijing),
	Subjects:      make(map[uuid.UUID][]ConvICS.Subject),
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("please appoint input file")
	}

	// open file
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatal(err)
	}

	err = parse(doc)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(schedule)
	log.Println(ioutil.WriteFile("schedule.ics", schedule.ToICS(timetable), 0777))
}

func parse(doc *goquery.Document) (err error) {
	defer func() {
		if Err := recover(); Err != nil {
			err = Err.(error)
			return
		}
	}()

	doc.Find("table").Each(func(i int, s *goquery.Selection) {
		switch i {
		case 0: //课表信息info
			parseTableInfo(s)
		case 1: //课表当前列
			parseSchedule(s)
		case 3:
			parseTimetable(s)
		}
	})
	return
}

func parseTimetable(s *goquery.Selection) {
	s.Find("tr").Each(func(i int, s *goquery.Selection) {
		if i < 3 { // 前面的不是我们校区
			return
		}
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			//if i == 0 { // 忽略校区名
			//	return
			//}
			var (
				desc   string
				sh, sm int
				eh, em int
			)
			_, err := fmt.Sscanf(s.Text(), "%s %d:%d~\n%d:%d", &desc, &sh, &sm, &eh, &em)
			if err != nil { //解析失败忽略
				log.Println("解析失败: ", err)
				return
			}

			start := time.Hour*time.Duration(sh) + time.Minute*time.Duration(sm)
			end := time.Hour*time.Duration(eh) + time.Minute*time.Duration(em)

			log.Printf("解析成功: %q 课: 从%d:%d到%d:%d", desc, sh, sm, eh, em)
			timetable.Append(desc, start, end) // 添加到Timetable
		})
	})
}

func parseTableInfo(s *goquery.Selection) {
	n := s.Find("tr > td")
	n.Each(func(i int, s *goquery.Selection) {
		switch i {
		case 1:
			log.Println("解析成功: 学期", s.Text())
		case 2:
			field := strings.Fields(s.Text())

			log.Println("解析成功: 学号: ", strings.TrimPrefix(field[0], "学号:"))
			log.Println("解析成功: 姓名: ", strings.TrimPrefix(field[1], "学生姓名:"))
			log.Println("解析成功: 班级: ", field[3])
			log.Println("解析成功: 学分: ", strings.TrimPrefix(field[4], "总学分:"))
		}
	})
}

var coursePat = regexp.MustCompile(`([^\s]+)\s+\(([^\s]+)\)\s+\(([^\s]+)\)\s+\(([^\s]+)\s+([^\s]+)\)`)

func parseSchedule(s *goquery.Selection) {
	s.Find("tr").EachWithBreak(func(n int, s *goquery.Selection) bool {
		s.Find("td").Each(func(j int, s *goquery.Selection) {
			if j == 0 { // 第一行是第几节课
				fmt.Printf("正在解析%q\n", strings.TrimSpace(s.Text()))
			} else {
				var duration int
				_, _ = fmt.Sscan(s.AttrOr("rowspan", "2"), &duration) //课长
				//单个课程占5行，把数据读取到结构体内
				//0
				//1 课程名 (编号)
				//2 (教师)
				//3
				//4 (n-m 地址)

				ans := coursePat.FindAllStringSubmatch(s.Text(), -1)
				for i := range ans {
					if len(ans[i]) != 5 {
						c := ans[i][1:]
						// c = ["大学物理实验I" "0268.46" "邓晓颖" "2-9" "物理实验室(未央)"]
						log.Printf("正在解析课程: [%d,%d]%s\n", n, j, c[0])
						findCourse(n, duration, time.Weekday(j), c[0], c[1], c[2], c[3], c[4])
					} else {
						log.Fatalf("解析失败: %q\n", ans[i])
					}
				}
			}
		})
		return true
	})
}

func findCourse(n, duration int, weekday time.Weekday, name, id, teacher, time, loc string) {
	log.Println(duration, name, id, teacher, time, loc)
	UUID := uuid.NewSHA1(uuid.NameSpaceX500, []byte(id))

	subject := schedule.Subjects[UUID]
	course := ConvICS.Subject{
		Name:     name,
		Teacher:  teacher,
		Location: loc,
		Weekday:  weekday,
		CStart:   n,
		CEnd:     duration,
	}
	if n, err := fmt.Sscanf(time, "%d-%d", &course.Start, &course.End); err != nil {
		if n == 1 {
			course.End = course.Start
		} else {
			log.Fatalf("解析周数%q出错: [%d]%v\n", time, n, err)
		}
	}

	subject = append(subject, course)
	schedule.Subjects[UUID] = subject
}
