package main

import "fmt"

func main() {
	map1 := make(map[string]string)
	map1["11"] = "11"

	fmt.Println("map是", map1)

	map2 := map1
	map1["22"] = "22"
	map2["33"] = "33"
	fmt.Println("map1是", map1)
	fmt.Println("map2是值拷贝?不是", map2)

	var slice1 []string
	slice1 = append(slice1, "aa")
	slice2 := slice1 //值拷贝，之后的操作相互不影响
	slice2 = append(slice2, "bb")
	fmt.Println("slice1", slice1)
	fmt.Println("slice2", slice2)

	var map3 = map[string]string{}
	map3["kk"] = ""
	fmt.Println(map3)
}
