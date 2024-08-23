package Draw

import (
	"ComicS/Model"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

// 输入的内容默认已经排序清洗了
// 只做一次下载
func SaveAllInfo(resplist []Model.SingeResult) ([]string, []string) {
	var picpathlist []string
	var txtpathlist []string

	currentTime := time.Now()
	timeStr := currentTime.Format("20060102_150405")
	dirpath := "result" + timeStr

	err := os.Mkdir(dirpath, 0755) // 0755 是目录的权限设置
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return nil, nil
	}

	for i, res := range resplist {
		//保存txt及封面
		newdirpath := dirpath + "/" + "数据" + strconv.Itoa(i)
		err := os.MkdirAll(newdirpath, 0755) // 0755 是目录的权限设置
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return nil, nil
		}
		txtpath := newdirpath + "/" + "展子讯息.txt"
		picpath := newdirpath + "/" + "cover.jpeg"
		datatemp := res.Conv2Com()
		writetxt(datatemp.String(), txtpath)
		downloadpic(datatemp.Cover, picpath)
		picpathlist = append(picpathlist, picpath)
		txtpathlist = append(txtpathlist, txtpath)
	}

	return picpathlist, txtpathlist
}

func Piclong(picpathlist []string) {
	// 存储所有处理后的图片
	var processedImages []image.Image

	// 目标大小
	targetWidth := 600
	targetHeight := 800

	// 读取并处理每张图片
	for _, file := range picpathlist {
		img, err := imaging.Open(file)
		if err != nil {
			log.Fatalf("failed to open image: %v", err)
		}

		// 将图片缩放到600x800像素
		resizedImg := imaging.Fit(img, targetWidth, targetHeight, imaging.Lanczos)

		// 将处理后的图片添加到列表中
		processedImages = append(processedImages, resizedImg)
	}

	// 计算输出图像的总尺寸
	outputWidth := targetWidth
	outputHeight := targetHeight * len(processedImages)

	// 创建一个新的RGBA图像作为最终拼接结果
	outputImage := imaging.New(outputWidth, outputHeight, color.NRGBA{255, 255, 255, 255})

	// 将每张图片按顺序绘制到输出图像上
	yOffset := 0
	for _, img := range processedImages {
		rect := image.Rect(0, yOffset, targetWidth, yOffset+targetHeight)
		draw.Draw(outputImage, rect, img, image.Point{0, 0}, draw.Over)
		yOffset += targetHeight
	}

	// 保存拼接后的图像
	currentTime := time.Now()
	timeStr := currentTime.Format("20060102_150405")
	err := imaging.Save(outputImage, timeStr+"long.jpg")
	if err != nil {
		log.Fatalf("failed to save output image: %v", err)
	}

	return
}

func PicSingle(picpathlist []string, resplist []Model.SingeResult) {
	// 目标大小
	targetWidth := 600
	targetHeight := 800

	outputWidth := targetWidth + 400
	outputHeight := targetHeight + 100

	currentTime := time.Now()
	timeStr := currentTime.Format("20060102_150405")
	dirpath := "./resultpic/" + timeStr

	err := os.MkdirAll(dirpath, 0755) // 0755 是目录的权限设置
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	for i, picpath := range picpathlist {
		img, _ := imaging.Open(picpath)
		resizedImg := imaging.Fit(img, targetWidth, targetHeight, imaging.Lanczos)
		outputImage := imaging.New(outputWidth, outputHeight, color.NRGBA{0, 0, 0, 0})
		rect := image.Rect(50, 50, targetWidth, 50+targetHeight)
		draw.Draw(outputImage, rect, resizedImg, image.Point{0, 0}, draw.Over)
		outputImagepath := dirpath + "/" + strconv.Itoa(i) + ".jpg"

		dc := gg.NewContextForImage(outputImage)
		// 设置矩形的颜色和透明度
		dc.SetColor(color.NRGBA{R: 128, G: 128, B: 128, A: 128})
		// 绘制一个矩形 (x, y, width, height)
		x := 50.0 + float64(targetWidth) + 25.0
		y := 50.0
		width, height := 300.0, 100.0
		dc.DrawRectangle(x, y, width, height)
		dc.Fill()
		// 设置文本的颜色
		dc.SetColor(color.White)
		// 设置字体大小
		fontSize := 20.0
		// 加载字体文件（你需要确保有一个 .ttf 文件）
		err = dc.LoadFontFace("./ttf/微软雅黑粗体.ttf", fontSize)
		if err != nil {
			log.Fatalf("failed to load font: %v", err)
		}

		// 绘制文本
		text1 := resplist[i].ProjectName
		// 绘制自动换行的文本
		drawWrappedString(dc, text1, x+10, y+10, width-20, fontSize)

		// 将绘制的内容保存回图像
		outputImg := dc.Image()
		if err := imaging.Save(outputImg, outputImagepath); err != nil {
			log.Fatalf("failed to save image: %v", err)
		}
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

// drawWrappedString 在指定区域内绘制自动换行的字符串
func drawWrappedString(dc *gg.Context, s string, x, y, maxWidth, lineHeight float64) {
	words := strings.Split(s, "") // 将中文文本按字符分割
	line := ""
	for _, word := range words {
		testLine := line + word
		w, _ := dc.MeasureString(testLine)
		if w > maxWidth {
			dc.DrawStringAnchored(line, x, y, 0, 0)
			y += lineHeight
			line = word
		} else {
			line = testLine
		}
	}
	if line != "" {
		dc.DrawStringAnchored(line, x, y, 0, 0)
	}
}
