# ConvICS
把教务网导出的课表转换为日历iCalendar文件。
也就是实现课表 -> .ics -> “日历”APP的转换。

每个学校教务网导出的格式都非常不同，所以你需要编写代码来将其转换为一个Schedule对象，
然后调用`func (s Schedule) ToICS(t Timetable) []byte`即可生成可被各大日历软件导入的ics文件。

## 已完成
> Feel free to P.R. your school's parser.

学校 | 方式 | 链接
-|-|-
西安工业大学 | 转换从教务网导出的xls文件 | [xatu](https://github.com/Tnze/ConvICS/tree/master/xatu) |
