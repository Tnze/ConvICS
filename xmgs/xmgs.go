// 解析教务网上导出的课程表
package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/Tnze/ConvICS"
	"log"
	"os"
	"strings"
	"time"
)

var timetable ConvICS.Timetable

/*
new(ConvICS.Timetable).
Append("第一节", "8:00AM", "8:50AM").
Append("第二节", "9:00AM", "9:50AM").
Append("第三节", "10:10AM", "11:00AM").
Append("第四节", "11:10AM", "12:00PM").
Append("第五节", "2:00PM", "2:50PM").
Append("第六节", "3:00PM", "3:50PM").
Append("第七节", "4:10PM", "5:00PM").
Append("第八节", "5:10PM", "6:00PM").
Append("第九节", "6:10PM", "7:20PM"). // 晚饭时间，一般不排课
Append("第十节", "7:30PM", "8:20PM").
Append("第十一节", "8:30PM", "9:20PM").
Append("第十二节", "9:30PM", "10:20PM")
*/

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
			if i == 0 { // 忽略校区名
				return
			}
			var (
				desc   string
				sh, sm int
				eh, em int
			)
			_, err := fmt.Sscanf(s.Text(), "%s %d:%d~\n%d:%d", &desc, &sh, &sm, &eh, &em)
			if err != nil { //解析失败忽略
				return
			}

			start := time.Date(0, 0, 0, sh, sm, 0, 0, time.UTC)
			end := time.Date(0, 0, 0, eh, em, 0, 0, time.UTC)

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

func parseSchedule(s *goquery.Selection) {
	s.Find("tr").EachWithBreak(func(n int, s *goquery.Selection) bool {
		s.Find("td").Each(func(j int, s *goquery.Selection) {
			if j == 0 { // 第一行是第几节课
				fmt.Printf("正在解析%q\n", strings.TrimSpace(s.Text()))
			} else {
				log.Println(s.Text())
				fe := strings.Split(s.Text(), "\n")

				for i := 0; i+5 < len(fe); i += 5 {
					var (
						duration int
						name     string
						teacher  string
					)
					fmt.Sscan(s.AttrOr("rowspan", "2"), &duration) //课长
					//单个课程占5行，把数据读取到结构体内
					//0
					//1 课程名 (编号)
					//2 (教师)
					//3
					//4 (n-m 地址)
					name = strings.TrimSpace(fe[i+1])
					teacher = strings.Trim(fe[i+2], "( )")
					//plan := strings.Trim(fe[i+4], "( )") + ")"
					//if _, err := fmt.Sscanf(plan, "%d-%d %s", &c.N, &c.M, &c.Loc); err != nil {
					//	fmt.Sscanf(plan, "%d %s", &c.N, &c.Loc)
					//}

					log.Println("课程: ", name, teacher, duration)
				}
			}
		})
		return true
	})
}
