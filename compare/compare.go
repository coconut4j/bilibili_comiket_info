package compare

import (
	"ComicS/Model"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func CRes(oldMap, curMap map[int]Model.SingeResult) (add, up []int) {

	for k, v := range curMap {
		v1, ok := oldMap[k]
		if !ok {
			//说明是新增的项目
			add = append(add, k)
			continue
		}
		//之前就有 比较
		equal := reflect.DeepEqual(v, v1)
		if !equal {
			//有修改 放到结果集里面
			up = append(up, k)
			continue
		}
	}

	return
}

func GetAddDetail(add []int, curMap map[int]Model.SingeResult) (addC map[int]*Model.CompareResult) {
	for _, k := range add {
		v := curMap[k]
		com := v.Conv2Com()
		addC[k] = com
	}
	return
}

// 先old
func GetUpDetail(up []int, oldMap, curMap map[int]Model.SingeResult) (upC map[int][2]*Model.CompareResult) {
	for _, k := range up {
		v1 := oldMap[k]
		oldCom := v1.Conv2Com()
		v2 := curMap[k]
		curCom := v2.Conv2Com()

		upC[k] = [2]*Model.CompareResult{oldCom, curCom}
	}
	return
}

func Format(add map[int]*Model.CompareResult, up map[int][2]*Model.CompareResult) string {
	var s strings.Builder
	if len(add) != 0 {
		for _, v := range add {
			s.WriteString(v.String())
			s.WriteString("---分割线---\r\n")
			s.WriteString("\r\n")
		}

	}
	if len(up) != 0 {
		s.WriteString("漫展变更 \r\n")
		marshal2, err := json.Marshal(up)
		if err != nil {
			fmt.Println(marshal2)
		}
		s.WriteString(string(marshal2))
	}

	return s.String()
}

func FormatOnTime(add map[int]*Model.CompareResult) string {
	var s strings.Builder
	if len(add) != 0 {
		s.WriteString("现有漫展数据 \r\n")
		s.WriteString("---分割线---\r\n")
		for _, v := range add {
			s.WriteString(v.String())
			s.WriteString("---分割线---\r\n")
			s.WriteString("\r\n")
		}

	} else {
		s.WriteString("未查询到深圳漫展数据")
	}
	return s.String()
}
