package main

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

func download(imageUrl, savePath, resolution string)  {
	if imageUrl == ""{
		fmt.Println("下载图片地址不能为空")
		return
	}
	if savePath == "" {
		fmt.Println("存储路径不能为空")
		os.Exit(0);
	}

	imageName := fmt.Sprintf("%s___%s", resolution, path.Base(imageUrl))

	fmt.Printf("开始下载，图片地址为：%s, 分辨率为: %s \n", imageUrl, resolution)

	image, err := http.Get(imageUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer image.Body.Close()

	reader := bufio.NewReaderSize(image.Body, 32 * 1024)

	file, err := os.Create(savePath  + imageName);
	if err != nil{
		panic(err)
	}

	writer := bufio.NewWriter(file)

	written, _ := io.Copy(writer, reader)

	fmt.Printf("Download success File Total size: %d \n", written)

}

func PathExists(path string){
	_, err := os.Stat(path);
	if err != nil {
		fmt.Printf("输入存储地址不存在，自动创建, %s \n", path)
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Printf("文件夹创建失败 %s \n", err)
			os.Exit(0)
		}
	}
}

func GetWallhavenPage(page int, savePath string)  {
	// Request the HTML page.

	httpUrl := fmt.Sprintf("https://wallhaven.cc/toplist?page=%d", page)

	res, err := http.Get(httpUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(".preview").Each(func(i int, s *goquery.Selection) {
		val, boolVal := s.Attr("href")
		if boolVal {
			valRes, err := http.Get(val)
			if err != nil{
				log.Fatal(err)
			}
			defer valRes.Body.Close()
			if res.StatusCode != 200 {
				log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
			}
			valDoc, err := goquery.NewDocumentFromReader(valRes.Body)
			if err != nil {
				log.Fatal(err)
			}

			imagePath, boolVal :=  valDoc.Find("#wallpaper").Attr("src")
			if boolVal {
				download(imagePath, savePath, strings.Replace(valDoc.Find(".showcase-resolution").Text(), " ", "", -1))
			}
		}
	})
}

func main()  {
	args := os.Args

	if args == nil {
		fmt.Println("请填写需要下载的页数以及储存地址  wallhaven [page] [savePath]")
		return
	}

	page, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	if page <= 0 {
		fmt.Println("下载页数必须大于0")
		return
	}


	savePath := args[2]

	if len(savePath) <= 0 {
		fmt.Println("请填写存储地址")
		return
	}

	PathExists(savePath)

	for i := 1; i<= page; i++ {
		GetWallhavenPage(i, savePath)
	}


	fmt.Println("全部下载结束")
}
