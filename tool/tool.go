package tool

import (
	"ComicS/Model"
)

func Conv(src []Model.SingeResult) map[int]Model.SingeResult {
	resp := make(map[int]Model.SingeResult, 32)
	for _, result := range src {
		resp[result.ProjectId] = result
	}
	return resp
}
