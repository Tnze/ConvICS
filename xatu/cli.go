// +build !wasm

package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	if len(os.Args) < 2 {
		log.Print("please appoint input file")
		log.Fatalf("Usage: %s 12345678903_李小龙（学生课表）.xls", os.Args[0])
	}

	// open file
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatal(err)
	}

	info, timetable, schedule, err := parse(doc)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%#v", info)

	err = ioutil.WriteFile("schedule.ics", schedule.ToICS(timetable), 0777)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("create .ics file success")
}
