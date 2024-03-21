package main

import (
	"ComicS/Model"
	"ComicS/bili"
	"fmt"
)

var (
	oldMap     map[int]Model.SingeResult
	convOldMap = make(map[int]*Model.CompareResult, 30)
	ontimeS    string
	pushString string
)

func main() {
	show_info_local()
}

func show_info_local() {
	bc := bili.Get_client()

	res, err := bc.GetAllResult()
	if err != nil {
		panic(err)
	}

	fmt.Println(res[1].ProjectName)

	bc.Show_result(res)
}
