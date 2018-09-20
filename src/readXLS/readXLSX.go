package readXLS

import (
	"fmt"
	//"github.com/satori/go.uuid"
	//"gopkg.in/yaml.v2"
	"github.com/tealeg/xlsx"
	"os"
	//"time"
	"strings"
)

//DoXLSX 把xlsx文件转换为yml格式的课程
func DoXLSX(fileName string) {

	file, err := xlsx.OpenFile(fileName)
	if err != nil {
		panic(fmt.Errorf("读文件%s失败: %v", fileName, err))
	}

	var yml string
	yml += "classes: \n"

	for _, sheet := range file.Sheets {
		fmt.Println("正在打开表：" + sheet.Name)
		fmt.Println(sheet.Cell(2, 0).String())
		for i := 4; i < 15; i++ {
			for j := 1; j < 8; j++ {
				text := sheet.Cell(i, j).String()
				if text != "" {
					t := strings.Split(text, "\n") //上下两行分开

					f := strings.Split(t[0], " ") //第一行
					//fmt.Println("课程名: " + f[0])       //课程名
					//fmt.Println("课程编号: " + cut(f[1])) //课程编号
					//fmt.Println("教师: ", cut(f[2]))    //教师名字
					timeLoc := cut(t[1])
					timeLocs := strings.Split(timeLoc, " ")
					times := strings.Split(timeLocs[0], "-")
					yml += "  - " + "\n"
					yml += "    time: "
					switch i {
					case 4:
						yml += "一二"
					case 6:
						yml += "三四"
					case 8:
						yml += "五六"
					case 10:
						yml += "七八"
					case 12:
						yml += "九十十一"
					}
					yml += "\n"
					yml += fmt.Sprintf("    day: %d", j) + "\n"
					yml += "    name: " + f[0] + "\n"
					yml += "    teacher: " + cut(f[2]) + "\n"
					yml += "    week: " + "\n"
					yml += "      - " + times[0] + "\n"
					yml += "      - " + times[1] + "\n"
					yml += "    location: " + timeLocs[1] + "\n"
					yml += "    code: " + cut(f[1]) + "\n" //课程编号
				}
			}
		}
	}
	yml += `schedule: 
  早读: #课名
    - 0720:00+08
    - 0750:00+08
  一: 
    - 080000+08
    - 085000+08
  二: 
    - 090000+08
    - 095000+08
  一二: 
    - 080000+08 
    - 095000+08 
  三: 
    - 101000+08
    - 110000+08
  四: 
    - 111000+08
    - 120000+08
  三四: 
    - 101000+08
    - 120000+08
  五: 
    - 140000+08
    - 145000+08
  六: 
    - 150000+08
    - 155000+08
  五六: 
    - 140000+08
    - 155000+08
  七: 
    - 161000+08
    - 170000+08
  八: 
    - 171000+08
    - 180000+08
  七八: 
    - 161000+08 
    - 180000+08
  九: 
    - 193000+08
    - 202000+08
  十:
    - 203000+08
    - 212000+08
  九十十一:
    - 193000+08
    - 212000+08
  晚修: 
    - 213000+08
    - 222000+08
firstWeek: 20180902`
	ymlfile, err := os.OpenFile("output.yml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	defer ymlfile.Close()
	if err != nil {
		panic(err)
	}
	ymlfile.WriteString(yml)
	//fmt.Println(yml)
}

//去掉字符串第一个和最后一个字符
func cut(s string) string {
	sr := []rune(s)
	return string(sr[1 : len(sr)-1])
}
