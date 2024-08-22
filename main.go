package main

import (
	"ComicS/Model"
	"ComicS/bili"
	"flag"
	"fmt"
	"os"
	"strings"
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

func main() {
	args := os.Args

	area := flag.String("area", "gz", "place")
	mode := flag.String("mode", "local", "mode")
	isdebug := flag.Bool("d", false, "debug")
	clear := flag.Bool("s", false, "only show saling")
	flag.Parse()

	sreacharea = *area
	onlyshowsale = *clear

	if !*isdebug {
		debuglv = 0
	}

	if debuglv == 1 {
		ShowArgs(args)
	}

	if strings.ToLower(*mode) == "local" {
		Show_info_local()
	} else if strings.ToLower(*mode) == "file" {
		SaveFile()
	} else if strings.ToLower(*mode) == "daemon" {
		DaemonRun()
	}
}

func Show_info_local() {
	//调取b站api
	bc := bili.Get_client(sreacharea, onlyshowsale)

	//获取结果
	res, err := bc.GetAllResult()
	if err != nil {
		panic(err)
	}

	//打印到终端
	if debuglv == 1 {
		res = bc.SortByTime(res)
		bc.Show_result(res)
	}

}

func SaveFile() {
	bc := bili.Get_client(sreacharea, onlyshowsale)
	res, err := bc.GetAllResult()
	if err != nil {
		panic(err)
	}
	res = bc.SortByTime(res)
	bc.Save2file(res)
}

func DaemonRun() {
	bc := bili.Get_client(sreacharea, onlyshowsale)
	bc.DaemonMode()
}

func ShowArgs(input []string) {
	fmt.Println("Program name:", input[0])
	if len(input) > 1 {
		for i, arg := range input[1:] {
			fmt.Printf("参数 %d: %s\n", i+1, arg)
		}
	} else {
		fmt.Println("无参数.")
	}
}
