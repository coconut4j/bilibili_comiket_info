package main

import (
	"ComicS/Model"
	"ComicS/bili"
	"ComicS/compare"
	"ComicS/tool"
	"fmt"
	"github.com/robfig/cron/v3"
	"net/http"
)

var (
	oldMap     map[int]Model.SingeResult
	convOldMap = make(map[int]*Model.CompareResult, 30)
	ontimeS    string
	pushString string
)

func main() {

	var c = cron.New()
	c.AddFunc("@every30m", updateInfo)
	c.Start()
	go func() {
		updateInfo()
	}()

	http.HandleFunc("/ontime", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(ontimeS))
	})

	http.ListenAndServe(":8091", nil)
	//fmt.Println(pushString)

	//接收退出信号量 退出程序
}

func updateInfo() {
	bc := bili.GetBiliBiliClient()
	info, err := bc.Info()
	if err != nil {
		panic(err)
	}
	fmt.Println(info.String())

	result, err := bc.GetAllResult()
	if err != nil {
		panic(err)
	}

	curMap := tool.Conv(result)

	if len(oldMap) == 0 {
		oldMap = curMap
		for k, v := range oldMap {
			com := v.Conv2Com()
			convOldMap[k] = com
		}
		ontimeS = compare.FormatOnTime(convOldMap)
	} else {
		add, up := compare.CRes(oldMap, curMap)
		addC := compare.GetAddDetail(add, curMap)
		upC := compare.GetUpDetail(up, oldMap, curMap)
		pushString = compare.Format(addC, upC)
	}

}
