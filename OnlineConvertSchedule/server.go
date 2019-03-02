package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("00000000000_XXX(学生课表).xls.exe")
	if err != nil {
		panic(err)
	}
	var data fullTable
	data.parse(f)
	j, err := json.Marshal(data)
	fmt.Println(string(j), err)
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
							c.Day = j + 1
							fmt.Sscan(s.AttrOr("rowspan", "2"), &c.Duration) //课长
							//单个课程占5行，把数据读取到结构体内
							//0
							//1 课程名 (编号)
							//2 (教师)
							//3
							//4 (n-m 地址)
							c.Name = strings.TrimSpace(fe[i+1])
							c.Teacher = strings.Trim(fe[i+2], "( )")
							fmt.Sscanf(strings.Trim(fe[i+4], "( )")+")", "%d-%d %s", &c.N, &c.M, &c.Loc)

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
