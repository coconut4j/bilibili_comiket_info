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

func main() {
	args := os.Args

	area := flag.String("area", "gz", "place")
	mode := flag.String("mode", "local", "mode")
	isdebug := flag.Bool("d", false, "debug")
	flag.Parse()

	sreacharea = *area

	if !*isdebug {
		debuglv = 0
	}

	if debuglv == 1 {
		ShowArgs(args)
	}

	if strings.ToLower(*mode) == "local" {
		Show_info_local()
	}
}

func Show_info_local() {
	//调取b站api
	bc := bili.Get_client(sreacharea)

	//获取结果
	res, err := bc.GetAllResult()
	if err != nil {
		panic(err)
	}

	//打印到终端
	if debuglv == 1 {
		bc.Show_result(res)
	}

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
