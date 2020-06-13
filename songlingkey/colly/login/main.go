package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
	"strconv"
)
type Book struct{
	N_global_guid string
	N_lang_id string
	N_format_id string
	V_title string

}
type BatchBook struct{
	Data []Book
	Status int
	Count int

}

func main(){


	//start:=11
	//end:=start+50
	//for true{
	//	go thread(start,end)
	//	time.Sleep(2 * time.Second) //等五秒重试
    //    if(end==163){
    //    	break
	//	}
	//	start=end+1
	//	end=start+30
	//	if(end>163){
	//		end=163
	//	}
	//}
	go thread(11,40)



	forever := make(chan bool)
	<-forever

}

func thread1(start int,end int) {
	fmt.Println("开始进行", start, end)
}

func thread(start int,end int) {
	fmt.Println("开始进行",start,end)

	var sum int
	var curdir string
	var whileflag bool
	// create a new collector
	c := colly.NewCollector()

	//c := colly.NewCollector(
	//	colly.Async(true),
	//)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 10})
	// authenticate
	err := c.Post("http://sy.sinocomic.com/Home/Login/check.shtml", map[string]string{"username": "bzj", "password": "123456"})
	if err != nil {
		log.Fatal(err)
	}


	// attach callbacks after login
	c.OnResponse(func(r *colly.Response) {
		//log.Println("response received", string(r.Body))

		//var buffer bytes.Buffer
		//json.HTMLEscape(&buffer,r.Body)
		bookCollector:=c.Clone()
		bookCollector.OnResponse(func(r *colly.Response) {
				fmt.Println("response received","内容是空？",len(r.Body))
				if(len(r.Body)==0){
					whileflag=false
					return
				}
			   fmt.Println("转换成图片")
				r.Save(curdir+"//"+strconv.Itoa(sum)+".png")

			})
		b:=&BatchBook{}
		json.Unmarshal(r.Body,b)
		fmt.Println("b",b)
		//bb=b
		for i,v:=range b.Data{
			fmt.Println("打印", i,v.N_global_guid,v.V_title)
			curdir="f://newbook11-40//"+v.V_title
			os.MkdirAll(curdir,os.ModePerm)
			whileflag=true
			sum=1
			for whileflag {
				//url:="http://sy.sinocomic.com/Home/Encrypt/decryptMprBookJpeg/id/"+v.N_global_guid+"/page/"+strconv.Itoa(sum)+"/lang/1.shtml"
				url:="http://sy.sinocomic.com/Home/Encrypt/decryptBookJpeg/id/"+v.N_global_guid+"/page/"+strconv.Itoa(sum)+".shtml"
				bookCollector.Visit(url)
				fmt.Println("第",sum,url)
				sum++
			}

		}


	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnHTML("data", func(e *colly.HTMLElement) {
		fmt.Println("打印", e.Text)

	})

  //for i:=1;i<=163;i++{
  	for i:=start;i<=end;i++{
		  c.Visit("http://sy.sinocomic.com/Home/Index/ajaxGetComicBook?curr="+strconv.Itoa(i)+"&id=0&isPhone=false")
	  }
	c.Wait()

	//bookCollector.OnResponse(func(r *colly.Response) {
	//	fmt.Println("response received","内容")
	//	if(len(r.Body)==0){
	//		whileflag=false
	//		return
	//	}
    //    fmt.Println("转换成图片")
	//	r.Save(curdir+"//"+strconv.Itoa(sum)+".png")
	//
	//})

	//bookCollector.Visit("http://sy.sinocomic.com/Home/Encrypt/decryptMprBookJpeg/id/2/page/200/lang/1.shtml")
	//for i,v:=range bb.Data{
	//	fmt.Println("打印", i,v.N_global_guid,v.V_title)
	//	curdir="f://book//"+v.V_title
	//	os.MkdirAll(curdir,os.ModePerm)
	//	whileflag=true
	//	sum=1
	//	for whileflag {
	//		url:="http://sy.sinocomic.com/Home/Encrypt/decryptMprBookJpeg/id/"+v.N_global_guid+"/page/"+strconv.Itoa(sum)+"/lang/1.shtml"
	//		bookCollector.Visit(url)
	//		fmt.Println("第",sum,url)
	//		sum++
	//	}
	//
	//}



	// Instantiate default collector
	//c := colly.NewCollector(
	//	// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
	//	colly.AllowedDomains("hackerspaces.org", "wiki.hackerspaces.org"),
	//)
	//
	//// On every a element which has href attribute call callback
	//c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	//	link := e.Attr("href")
	//	// Print link
	//	fmt.Printf("Link found: %q -> %s\n", e.Text, link)
	//	// Visit link found on page
	//	// Only those links are visited which are in AllowedDomains
	//	c.Visit(e.Request.AbsoluteURL(link))
	//})
	//
	//// Before making a request print "Visiting ..."
	//c.OnRequest(func(r *colly.Request) {
	//	fmt.Println("Visiting", r.URL.String())
	//})
	//
	//// Start scraping on https://hackerspaces.org
	//c.Visit("https://hackerspaces.org/")
}