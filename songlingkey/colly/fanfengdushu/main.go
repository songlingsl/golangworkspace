package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	title := ""
	suffix := "mp3"
	curdir := ""
	sum := 200002437
	//var curdir string
	//var whileflag bool
	c := colly.NewCollector()
	down := colly.NewCollector(
		colly.MaxBodySize(0),
	)
	err := c.Post("https://api.dushu.io/fragment/content", map[string]string{"username": "13810938737", "password": "sl796800"})
	if err != nil {
		log.Fatal(err)
	}

	c.OnResponse(func(r *colly.Response) {

		log.Println(sum, "内容是：", string(r.Body))
		m := make(map[string]interface{})
		json.Unmarshal(r.Body, &m)
		if m["title"] == nil {
			return
		}
		title = m["title"].(string)
		type1 := m["type"].(float64)
		if type1 != 2 {
			return
		}
		mediaUrls := m["mediaUrls"].([]interface{})
		//m["mediaUrls"].
		if len(mediaUrls) <= 0 {
			return
		}
		curdir = "G://fanden1//" + title

		if strings.Contains(title, "曾国藩") {
			os.MkdirAll(curdir, os.ModePerm)
			r.Save(curdir + "//" + strconv.Itoa(sum) + ".txt")
			fmt.Println("名：", title)
			url := mediaUrls[0].(string)
			start := strings.LastIndex(url, ".")
			suffix = url[start:]
			down.Visit(url)
		}

	})
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Content-Type", "application/x-www-form-urlencoded")
	})

	c.OnHTML("", func(e *colly.HTMLElement) {
		fmt.Println("打印", e.Text)

	})

	down.OnResponse(func(r *colly.Response) {
		log.Println("视频或音频", curdir+"//"+strconv.Itoa(sum)+suffix)
		r.Save(curdir + "//" + strconv.Itoa(sum) + suffix)
	})
	down.OnRequest(func(r *colly.Request) {
		fmt.Println("downVisiting", r.URL.String())
	})

	//for i:=start;i<=end;i++{
	//	c.Visit("http://sy.sinocomic.com/Home/Index/ajaxGetComicBook?curr="+strconv.Itoa(i)+"&id=0&isPhone=false")
	//}
	//c.Visit("https://api.dushu.io/fragment/content?fragmentId=200002437&token=8CAdxcg0uP6eAY4f67Z")

	//for sum = 200001000; sum <= 200005000; sum++ {
	for sum = 1; sum <= 200001000; sum++ {
		fmt.Println("打印sum1", sum)
		c.Post("https://api.dushu.io/fragment/content", generateFormData(sum))
	}
	//}
	c.Wait()
	forever := make(chan bool)
	<-forever
}

func generateFormData(i int) map[string]string {
	return map[string]string{
		"fragmentId": strconv.Itoa(i),
		"token":      "8CAdxcg0uP6eAY4f67Z",
	}
}
