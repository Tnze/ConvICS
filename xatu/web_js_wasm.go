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
	info, timetable, schedule, err := parse(doc)
	if err != nil {
		log.Fatal(err)
	}

	// 生成ics
	ics := schedule.ToICS(timetable)
	if len(args) >= 2 && args[1].Type() == js.TypeFunction {
		data := js.Global().Get("Uint8Array").New(js.ValueOf(len(ics)))
		js.CopyBytesToJS(data, ics)
		args[1].Invoke(data)
	} else {
		fmt.Print("ignored the ics")
	}

	// 返回课表信息
	if len(args) >= 3 && args[2].Type() == js.TypeFunction {
		args[2].Invoke(map[string]interface{}{
			"year":  info.Year,
			"id":    info.ID,
			"name":  info.Name,
			"class": info.Class,
			"score": info.Score,
		})
	} else {
		fmt.Print("ignored the tableinfo")
	}

	return nil
}
