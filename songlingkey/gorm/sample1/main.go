package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Animal struct {//id自增
	AnimalId int64 `gorm:"PRIMARY_KEY:Animal_Id"`
	Name     string
	Age      int
}
func main() {
	db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3307)/msg3?charset=utf8&parseTime=True&loc=Local")
	if(err!=nil){
		fmt.Println("连不上",err.Error())
	}
	db.SingularTable(true)//设置后才能对应好表名
	defer db.Close()
	//a:=Animal{Name:"哈哈",Age:223}
	//db.Create(&a)//保存


	var animal Animal
	db.First(&animal, "Animal_Id = ?",24) // 查询Animal_Id为24的animal
	fmt.Println("得到：",animal.Name)


	animal2:=Animal{AnimalId:36}
	db.First(&animal2) // 查询Animal_Id为36的animal
	fmt.Println("得到：",animal2.Name)


	var list [ ]Animal//定义一个User类型的数组名字叫做users
	db.Where("name = ? AND age = ?", "哈哈", "223").Find(&list)//两个条件
	fmt.Println("得到：",list)

	db.Model(&animal).Where("Animal_Id = ?", 22).Updates(map[string]interface{}{"name": "我1", "age": 13})
	db.Model(&animal).Where("Animal_Id = ?", 23).Updates(Animal{Name: "hello", Age: 18})

	//db.Raw("SELECT name, age FROM users WHERE name = ?", 3).Scan(&result) 直接sql
}
