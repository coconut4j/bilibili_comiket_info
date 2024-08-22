package main

import (
	"ComicS/Model"
	"ComicS/bili"
	"flag"
	"fmt"
	"strings"
	"time"
)

var (
	oldMap     map[int]Model.SingeResult
	convOldMap = make(map[int]*Model.CompareResult, 30)
	ontimeS    string
	pushString string
)

var debuglv int = 1
var sreacharea string
var onlyshowsale bool
var unixendtime int64 = -1

func main() {
	area := flag.String("area", "gz", "place")
	mode := flag.String("mode", "local", "mode")
	isdebug := flag.Bool("d", false, "debug")
	clear := flag.Bool("s", false, "only show saling")
	endtimestring := flag.String("end", "no", "endtime format like 2024-08-23")
	flag.Parse()

	sreacharea = *area
	onlyshowsale = *clear

	if !*isdebug {
		debuglv = 0
	}

	if *endtimestring != "no" {
		t, _ := time.Parse("2006-01-02 15:04:05", *endtimestring)
		unixendtime = t.Unix()
		fmt.Printf("\n flag data %s unix time:%d \n", *endtimestring, unixendtime)
	}

	if strings.ToLower(*mode) == "local" {
		Show_info_local()
	} else if strings.ToLower(*mode) == "file" {
		SaveFile()
	} else if strings.ToLower(*mode) == "daemon" {
		DaemonRun()
	} else if strings.ToLower(*mode) == "pic" {
		genpic()
	}
}

func Show_info_local() {
	bc := bili.Get_client(sreacharea, onlyshowsale)
	res, err := bc.GetAllResult()
	if err != nil {
		panic(err)
	}
	if debuglv == 1 {
		res = bc.SortByTime(res, unixendtime)
		bc.Show_result(res)
	}

}

func SaveFile() {
	bc := bili.Get_client(sreacharea, onlyshowsale)
	res, err := bc.GetAllResult()
	if err != nil {
		panic(err)
	}
	res = bc.SortByTime(res, unixendtime)
	bc.Save2file(res)
}

func DaemonRun() {
	bc := bili.Get_client(sreacharea, onlyshowsale)
	bc.DaemonMode()
}

func genpic() {
	bc := bili.Get_client(sreacharea, onlyshowsale)
	res, err := bc.GetAllResult()
	if err != nil {
		panic(err)
	}
	res = bc.SortByTime(res, unixendtime)
	bc.Pic(res)
}
