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
	bilibiliHost    = "https://show.bilibili.com/"
	ShenzhenArea    = "440300"
	DefaultPageSize = "16"
	DefaultPage     = "1"
	DefaultType     = "展览"
)

var (
	c    *http.Client
	bc   *BiliClient
	once sync.Once
)

type BiliClient struct {
	c            *http.Client
	BilibiliHost string
}

func GetBiliBiliClient() *BiliClient {
	if bc == nil {
		once.Do(func() {
			c = &http.Client{Timeout: 10 * time.Second}
			bc = &BiliClient{
				c:            c,
				BilibiliHost: bilibiliHost,
			}
		})
	}
	return bc
}

func (bc *BiliClient) CreateUrl(page, pagesize, area, p_type string) string {
	s := fmt.Sprintf("%s/api/ticket/project/listV2?version=134&page=%s&pagesize=%s&area=%s&filter=&platform=web&p_type=%s", bc.BilibiliHost, page, pagesize, area, p_type)
	fmt.Println(s)
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

func (bc *BiliClient) Info() (Model.BaseInfo, error) {
	var resp Model.BaseInfo
	data, err := bc.GetFirstPageData()
	if err != nil {
		return resp, err
	}

	resp.Total = data.Total
	resp.Page = data.Page
	resp.Pagesize = data.Pagesize
	resp.NumResults = data.NumResults
	resp.NumPages = data.NumPages
	return resp, nil
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

/*
func a() {
	//计算总数和页数

	//从第一页开始

	firstPage := do(&c, url)
	allData = append(allData, firstPage)
	//总数
	//total := firstPage.Total
	//页数
	NumPages := firstPage.NumPages

	if NumPages == 1 {
		//todo 直接处理结果
	}
	//第一页已经请求过了 从第二页开始多线程请求 且存在第二页
	wg.Add(NumPages - 1) //需要numpages-1个协程等待
	respchan := make(chan *Data, NumPages)
	for i := 2; i <= NumPages; i++ {
		url = fmt.Sprintf("/api/ticket/project/listV2?version=134&page=%d&pagesize=16&area=440300&filter=&platform=web&p_type=展览", i)
		go func() {
			defer wg.Done()
			a := do(&c, url)

			respchan <- a
		}()
	}
	wg.Wait()
	close(respchan)
	for data := range respchan {
		allData = append(allData, data)
	}
	for page, datum := range allData {
		fmt.Printf("page = %d \r\n data = %v \r\n", page+1, datum)
	}
}

func do(c *http.Client, url string) {
	resp, err := c.Get(BilibiliHost + url)
	if err != nil {
		panic(err)
	}

}
*/
