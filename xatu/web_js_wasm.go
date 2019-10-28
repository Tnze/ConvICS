// +build wasm

package main

import (
	"bytes"
	"fmt"
	"log"
	"syscall/js"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	js.Global().Set("ConvToICS", js.FuncOf(convert))
	select {}
}

func convert(this js.Value, args []js.Value) interface{} {
	// 读取数据
	fmt.Println("Converting")
	data := make([]byte, args[0].Length())
	js.CopyBytesToGo(data, args[0])

	// 解析数据
	r := bytes.NewReader(data)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Fatal(err)
	}

	// 提取课表
	err = parse(doc)
	if err != nil {
		log.Fatal(err)
	}

	// 生成ics
	schedule.ToICS(timetable)

	return nil
}
