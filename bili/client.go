package bili

import (
	"ComicS/Model"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	Page_Size       = 1 << 4
	baseurl         = "https://show.bilibili.com"
	ShenzhenArea    = "440300"
	DefaultPageSize = "16"
	DefaultPage     = "1"
	DefaultType     = "展览"
)

var (
	c  *http.Client
	bc *BiliClient
)

type BiliClient struct {
	c            *http.Client
	BilibiliHost string
}

func Get_client() *BiliClient {
	if bc == nil {
		c = &http.Client{Timeout: 10 * time.Second}
		bc = &BiliClient{
			c:            c,
			BilibiliHost: baseurl,
		}
	}
	return bc
}

func (bc *BiliClient) CreateUrl(page, pagesize, area, p_type string) string {
	s := fmt.Sprintf("%s/api/ticket/project/listV2?version=134&page=%s&pagesize=%s&area=%s&filter=&platform=web&p_type=%s", bc.BilibiliHost, page, pagesize, area, p_type)
	return s
}

func (bc *BiliClient) GetDefaultUrl() string {
	return bc.CreateUrl(DefaultPage, DefaultPageSize, ShenzhenArea, DefaultType)
}

func (bc *BiliClient) GetOnePageData(page, pagesize string) (*Model.Data, error) {
	url := bc.CreateUrl(page, pagesize, ShenzhenArea, DefaultType)
	return bc.request(url)
}

func (bc *BiliClient) GetFirstPageData() (*Model.Data, error) {
	return bc.GetOnePageData("1", DefaultPageSize)
}

func (bc *BiliClient) request(url string) (*Model.Data, error) {
	resp, err := bc.c.Get(url)
	if err != nil {
		return nil, err
	}
	respStruct := &Model.RespAll{}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, respStruct)
	if err != nil {
		return nil, err
	}
	if respStruct.Errno != 0 {
		return nil, errors.New("http errmsg = " + respStruct.Msg)
	}
	return &respStruct.Data, err
}

func (bc *BiliClient) GetAllResult() ([]Model.SingeResult, error) {
	var resp []Model.SingeResult

	data, err := bc.GetFirstPageData()
	if err != nil {
		return nil, err
	}

	//只有一页数据直接返回
	numPages := data.NumPages
	if numPages == 1 {
		return data.Result, err
	}
	resp = append(resp, data.Result...)

	respchan := make(chan *Model.Data, numPages-1)
	errChan := make(chan error, numPages-1)
	wg := sync.WaitGroup{}
	wg.Add(numPages - 1)
	for i := 2; i <= numPages; i++ {
		pageIndex := strconv.FormatInt(int64(i), 10)
		go func(pageIndex string) {
			defer wg.Done()
			pageData, err := bc.GetOnePageData(pageIndex, DefaultPageSize)
			respchan <- pageData
			if err != nil {
				errChan <- err
			}
		}(pageIndex)
	}
	wg.Wait()
	close(errChan)
	close(respchan)
	for err2 := range errChan {
		if err2 != nil {
			fmt.Println("调用发生了错误")
			return nil, err2
		}
	}
	for datum := range respchan {
		resp = append(resp, datum.Result...)
	}
	return resp, nil
}

func (bc *BiliClient) Show_result(resplist []Model.SingeResult) {
	for _, res := range resplist {
		maptemp := res.Conv2Com()
		fmt.Println(maptemp.String())
	}

	return
}
