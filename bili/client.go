package bili

import (
	"ComicS/Model"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	Page_Size       = 1 << 4
	baseurl         = "https://show.bilibili.com"
	SzArea          = "440300"
	DefaultPageSize = "16"
	DefaultPage     = "1"
	DefaultType     = "展览"
	GzArea          = "440100"
)

var (
	c          *http.Client
	bc         *BiliClient
	sreacharea string
	isclear    bool
)

type BiliClient struct {
	c            *http.Client
	BilibiliHost string
}

func Get_client(area string, cleardata bool) *BiliClient {
	sreacharea = area
	isclear = cleardata
	if bc == nil {
		c = &http.Client{Timeout: 10 * time.Second}
		bc = &BiliClient{
			c:            c,
			BilibiliHost: baseurl,
		}
	}
	return bc
}

// b站默认是获取第一页后确定pagesize 默认1开始
func (bc *BiliClient) CreateUrl(page, pagesize, area, p_type string) string {
	s := fmt.Sprintf("%s/api/ticket/project/listV2?version=134&page=%s&pagesize=%s&area=%s&filter=&platform=web&p_type=%s", bc.BilibiliHost, page, pagesize, area, p_type)
	return s
}

func (bc *BiliClient) GetDefaultUrl() string {
	var temp string
	if strings.ToLower(sreacharea) == "gz" {
		temp = GzArea
	} else if strings.ToLower(sreacharea) == "sz" {
		temp = SzArea
	} else {
		temp = GzArea
	}
	return bc.CreateUrl(DefaultPage, DefaultPageSize, temp, DefaultType)
}

func (bc *BiliClient) GetOnePageData(page, pagesize string) (*Model.Data, error) {
	var temp string
	if strings.ToLower(sreacharea) == "gz" {
		temp = GzArea
	} else if strings.ToLower(sreacharea) == "sz" {
		temp = SzArea
	} else {
		temp = GzArea
	}
	url := bc.CreateUrl(page, pagesize, temp, DefaultType)
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
	count := 0
	for _, res := range resplist {
		if isclear && res.SaleFlagNumber != 2 {
			continue
		}
		maptemp := res.Conv2Com()
		fmt.Println(maptemp.String())
		count++
	}

	fmt.Printf("\n总计%d个有效信息 收集了%d个信息\n", count, len(resplist))

	return
}

func (bc *BiliClient) SortByTime(resplist []Model.SingeResult, unixendtime int64) []Model.SingeResult {
	if isclear {
		resplist = delunuse(resplist, unixendtime)
	}
	sort.Sort(Model.ByUnixTimestamp(resplist))
	return resplist
}

func delunuse(resplist []Model.SingeResult, unixendtime int64) []Model.SingeResult {
	var output []Model.SingeResult

	if unixendtime == -1 {
		for _, res := range resplist {
			if res.SaleFlagNumber == 2 {
				output = append(output, res)
			}
		}
	} else {
		for _, res := range resplist {
			if res.SaleFlagNumber == 2 && res.StartUnix <= unixendtime {
				output = append(output, res)
			}
		}
	}

	return output
}

func (bc *BiliClient) Save2file(resplist []Model.SingeResult) {
	//在当前目录下利用时间戳建立文件夹保存
	currentTime := time.Now()
	timeStr := currentTime.Format("20060102_150405")
	dirpath := "result" + timeStr

	err := os.Mkdir(dirpath, 0755) // 0755 是目录的权限设置
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	for i, res := range resplist {
		//保存txt及封面
		newdirpath := dirpath + "/" + "数据" + strconv.Itoa(i)
		err := os.MkdirAll(newdirpath, 0755) // 0755 是目录的权限设置
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
		txtpath := newdirpath + "/" + "展子讯息.txt"
		picpath := newdirpath + "/" + "cover.jpeg"
		datatemp := res.Conv2Com()
		writetxt(datatemp.String(), txtpath)
		downloadpic(datatemp.Cover, picpath)

	}

	return
}

func writetxt(txtdata string, savepath string) error {
	// 创建或打开文件
	file, err := os.Create(savepath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close() // 确保在函数结束时关闭文件

	// 写入文本到文件
	_, err = file.WriteString(txtdata)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}

	return nil
}

func downloadpic(url string, savepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 创建文件
	out, err := os.Create(savepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 将 HTTP 响应的主体（即图片数据）复制到文件中
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (bc *BiliClient) DaemonMode() {

	envset()

	if os.Getenv("GO_DAEMON") != "1" {
		cmd := exec.Command(os.Args[0])
		cmd.Env = append(os.Environ(), "GO_DAEMON=1")
		cmd.Stdout = nil
		cmd.Stderr = nil
		cmd.Stdin = nil

		sysinfo := runtime.GOOS
		if sysinfo != "windows" {
			//unix系统专供daemon win得服务注册 unix编译可以删除备注
			//cmd.SysProcAttr = &syscall.SysProcAttr{
			//	Setsid: true,
			//}
		}

		err := cmd.Start()
		if err != nil {
			log.Fatal("Failed to start daemon process: ", err)
		}
		log.Println("Daemon process started with PID:", cmd.Process.Pid)
		os.Exit(0)
	}

	rpcmode()

	return
}

func (bc *BiliClient) Pic(resplist []Model.SingeResult) {

	return
}

func envset() {

}

func rpcmode() {
	// 关闭标准输入、输出和错误流
	f, err := os.OpenFile("/dev/null", os.O_RDWR, 0)
	if err != nil {
		log.Fatal("Failed to open /dev/null: ", err)
	}
	defer f.Close()
	os.Stdout = f
	os.Stderr = f
	os.Stdin = f

	// 模拟守护进程的持续工作
	for {
		//预留为protoc grpc提供http查询业务
		//proto buffer编写数据及api格式 预留能力
	}
}
