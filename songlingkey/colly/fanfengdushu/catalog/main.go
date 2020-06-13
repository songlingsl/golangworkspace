package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)
type Content struct{
	FragmentId int
	Type int

}
type Book struct{
	BookId string `gorm:"PRIMARY_KEY:book_Id"`
	Title string `gorm:"-"`
	Summary string
	ImageUrl string
	Recommendation string
	Contents []Content `gorm:"-"`
	BookName string
	MediaUrl string
	ShowTime time.Time
	ReadCount  int
	CreateTime int64 `gorm:"-"`
}
type Jsonobj struct{
	Books  []Book


}
var db *gorm.DB
var dbbook Book
func main() {
	var err error
	db, err = gorm.Open("mysql", "root:root@tcp(127.0.0.1:3307)/mybook?charset=utf8mb4&parseTime=True&loc=Local")

	if(err!=nil){
		fmt.Println("连不上",err.Error())
	}
	db.SingularTable(true)//设置后才能对应好表名

	suffix:="mp3"
	curdir:=""
	title:=""
	FragmentId:=0
	c := colly.NewCollector()
	d := colly.NewCollector()
	down:= colly.NewCollector(
		colly.MaxBodySize(0),

	)
	err = d.Post("https://api.dushu.io/fragment/content", map[string]string{"username": "13810938737", "password": "sl796800"})
	if err != nil {
		log.Fatal(err)
	}
	c.OnResponse(func(r *colly.Response) {
		log.Println( "内容是：", string(r.Body))

		objs:=Jsonobj{}
		json.Unmarshal(r.Body,&objs)

		for _, book := range objs.Books {
			//fmt.Println("书名",book.Title)
			title=book.Title
			book.BookName=book.Title

			fmt.Println("时间啊",book)


			book.ShowTime=time.Unix(book.CreateTime/1000,0)
			dbbook=book
			curdir="G://fandengnew//"+book.Title
			fmt.Println("书：",book)
			os.MkdirAll(curdir,os.ModePerm)
			for _, c := range book.Contents{
				FragmentId=c.FragmentId
				//fmt.Println("片：",c.FragmentId)
				//不用保存mp3了
				d.Post("https://api.dushu.io/fragment/content", fragmentData(c.FragmentId))
			}

		}

	})
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Content-Type", "application/json; charset=UTF-8")
	})


	d.OnRequest(func(r *colly.Request) {

		r.Headers.Set("Content-Type", "application/x-www-form-urlencoded")
	})

	d.OnResponse(func(r *colly.Response) {

		log.Println(FragmentId,"内容是：", string(r.Body))
		m := make(map[string]interface{})
		json.Unmarshal(r.Body, &m)
		if( m["title"]==nil){
			return
		}
		title = m["title"] .(string)
		type1 := m["type"].(float64)
		if(type1!=2){
			return
		}
		mediaUrls := m["mediaUrls"] .( []interface{})
		//m["mediaUrls"].
		if(len( mediaUrls)<=0){
			return
		}
		//r.Save(curdir+"//"+strconv.Itoa(FragmentId)+".txt")

		//fmt.Println("书名：",title)
		mediaurl:=mediaUrls[0].(string)
		start:=strings.LastIndex(mediaurl,".")
		suffix=mediaurl[start:]
		dbbook.MediaUrl=mediaurl
		//if(strings.Contains(title,"巴菲特")){
		//	fmt.Println("书名有巴菲特：",url)
		//	down.Visit(url)
		//}
		//down.Visit(mediaurl)
		saveDB()

	})

	down.OnResponse(func(r *colly.Response) {
		log.Println(title,"的视频或音频",curdir+"//"+strconv.Itoa(FragmentId)+suffix)
		//r.Save(curdir+"//"+strconv.Itoa(FragmentId)+suffix)


	})

	down.SetRequestTimeout(120*time.Second)//必须设置下载的时间，不然下载不全或没有
	//for i:=0;i<8;i++{
	c.PostRaw("https://api.dushu.io/books", generateFormData(0))
	//}
	c.Wait()
	d.Wait()
	down.Wait()
	forever := make(chan bool)
	<-forever

}

func saveDB() {
	fmt.Println("保存的书：",dbbook)
	db.Create(&dbbook)//保存
}


func generateFormData(i int) []byte {
	str:="{\"pageSize\":1000,\"page\":1,\"categoryId\":"+strconv.Itoa(i)+",\"token\":\"8CAdxcg0uP6eAY4f67Z\",\"order\":1,\"bookReadStatus\":-1}"
   fmt.Println("参数：",str)
    return []byte(str)
}
func fragmentData(i int) map[string]string {
	return map[string] string{
		"fragmentId": strconv.Itoa(i),
		"token":  "8CAdxcg0uP6eAY4f67Z",
	}
}