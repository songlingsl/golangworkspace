package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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
	BookType string
	CreateTime int64 `gorm:"-"`
}
type Jsonobj struct{
	Books  []Book


}
var db *gorm.DB
func main() {

	var err error
	db, err = gorm.Open("mysql", "root:root@tcp(127.0.0.1:3307)/mybook?charset=utf8mb4&parseTime=True&loc=Local")

	if(err!=nil){
		fmt.Println("连不上",err.Error())
	}
	db.SingularTable(true)//设置后才能对应好表名

	file, err := os.Open("G:\\fandengtype\\作者光临.txt")   //打开
	if err != nil { fmt.Println(err); return  }
	defer file.Close() //关闭

	line := bufio.NewReader(file)
	for {
		content, _, err := line.ReadLine()
		if err == io.EOF { break }
		strc:=string(content)
		if strings.Contains(strc,"title"){
			c:=strings.Split(strc,":")
			value:=strings.ReplaceAll(string(c[1]),"\"","")
			value=strings.TrimSpace(value)
			updateDB(value,"作者光临")
		}

	}

}

func updateDB(value string,types string){
    fmt.Println("书名"+value+"类型"+types)
	db.Model(&Book{}).Where("book_name = ?", value).Updates(map[string]interface{}{"bookType": types})
}

