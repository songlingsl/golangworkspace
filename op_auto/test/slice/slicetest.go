package main

import "fmt"

func main() {
	//对于很多数据来讲：频繁的插入和删除用list,频繁的遍历查询选slice
	var array1 [2]string //数组
	array1[0] = "0"
	array1[1] = "1"
	//array1 = append(array1, "3")数组无法append,数量固定

	var slice1 []string         //切片
	slice2 := make([]string, 3) //切片，初始3个空元素，用append会自动扩充,以诚意2的数量扩容。最好不设置
	slice1 = append(slice1, "aaaaaaaaa")
	slice2 = append(slice2, "bb", "cc", "dd")
	slice2[0] = "ee"
	slice2 = append(slice2, "ff", "gg", "hh")
	//slice2[20] = "88" 越界，需要用append
	fmt.Println(len(slice2)) //非空数量
	fmt.Println(cap(slice2)) //容量
	fmt.Println(slice1, slice2)

	for v, i := range slice2 {

		fmt.Println(v, "卡", i)

	}

}
