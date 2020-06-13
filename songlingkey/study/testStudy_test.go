package main

import (
	"container/list"
	"context"
	"fmt"
	"os"
	"reflect"
	"runtime/debug"
	"runtime/pprof"
	"sync"
	"testing"
	"time"
)

//空接口  interface{}   ...interface{} java的Object
//强制转换  如果a是interface  转string  a.(string)
//golang 日志框架Logrus
func TestSlice(t *testing.T) {
	println("文件后缀是_test")
	slice1 := make([]string, 1)
	slice2 := append(slice1, "dd")
	println("值", slice2[1])

	slice2[0] = "00000"
	for i, v := range slice2 {
		println("222：", i, v)
	}

	for i, v := range slice1 {
		println("111：：", i, v)
	}
	var slice3 []string
	println(append(slice3, "xin"))

	slice := []int{10, 20, 30, 40, 50}
	// 创建一个新切片
	// 其长度为 2 个元素,容量为 4 个元素
	newSlice := slice
	newSlice[2] = 222
	// 使用原有的容量来分配一个新元素
	// 将新元素赋值为 60，会改变底层数组中的元素
	newSlice = append(newSlice, 60)
	fmt.Println(slice, newSlice)
	//slice2 = append(slice2, slice1...)//切片加切片 ...即使拆散

}

func TestSlice2(t *testing.T) {
	var slice1 [10]string //切片长度需要定义，不然超过长度，就会生成新切片(新底层数组),切片里都是nil也不好
	fmt.Println(len(slice1))

	var slice2 = make([]string, 0, 10) //预定义切片容量最好 append 不会生成新切片

	slice2 = append(slice2, "值1")
	slice3 := append(slice2, "值2")
	fmt.Println(slice2)
	fmt.Println(slice3)
	slice2[0] = "值3"
	fmt.Println(slice2)
	fmt.Println(slice3)
	println("地址", slice2)         //底层地址
	fmt.Printf("addr:%p", slice2) //底层地址
	println("地址", slice3)         //底层地址一致
	fmt.Printf("addr:%p", slice3)

}

func TestDict(t *testing.T) {
	map1 := make(map[string]string)
	map1["11"] = "11"
	println(map1)
	fmt.Println(map1)

	for i, v := range map1 {
		println("迭代", i, v)
	}
}
func TestUnless(t *testing.T) {
	println("开始")
	unless := make(chan bool)
	<-unless
	println("结束")
}

type Person struct {
	Name string
	age  int
}

func (p *Person) Abc(ss string) string {
	return "person内部方法" + ss
}
func TestStruct(t *testing.T) {
	var p1 Person
	fmt.Println(p1.Abc("加"))
	p3 := p1
	if p3 == p1 {
		//不等于，值 拷贝
	}
	p4 := &p1
	if *p4 == p1 {
		fmt.Println("这个等于")
	}
	p1.Name = "宋"
	fmt.Println(p3.Name, "类型", reflect.TypeOf(p1))
	printP(p1)
	printPerson(&p1)
	printP(p1) //生效
	p2 := &Person{Name: "宋"}
	fmt.Println("类型", reflect.TypeOf(p2))
}

func printP(p Person) {
	println(p.Name)
	fmt.Println("类型", reflect.TypeOf(p))
	p.Name = "新" //不生效，因为不是指针
}

func printPerson(p *Person) {
	println(p.Name)
	fmt.Println("类型", reflect.TypeOf(p))
	p.Name = "新" //不生效，因为不是指针
}

func TestException(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("下面异常抛出来了", err)
			return
		}
	}()
	Exception()
	fmt.Println("可以继续执行")
	v := Exception1()
	println("打印这个值吗", v)
}

func Exception() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("异常啦", err)
			debug.PrintStack()
			return
		}
	}()
	abc := make([]string, 1)
	hah := abc[10]
	fmt.Println(hah)
}
func Exception1() string {
	abc := make([]string, 1)
	hah := abc[10]
	fmt.Println(hah)
	return "返回值"
}

var Locker sync.Mutex
var wait sync.WaitGroup

func TestGorutine(t *testing.T) {
	wait.Add(1)
	go Goru()
	wait.Wait()
	//Locker.Lock()
	//println("线程锁")
	//Locker.Unlock()
	println("解锁")
	unless := make(chan bool)
	<-unless
}
func Goru() {
	time.Sleep(5 * time.Second)
	println("停了几秒")
	wait.Done()

}

//同步型map
func TestSyncMap(t *testing.T) {
	var map1 sync.Map
	map1.Store("11", "22")
	fmt.Println(map1.Load("11"))
}
func TestList(t *testing.T) {
	mylist := list.New()
	mylist.PushBack("拉了")
	fmt.Println(mylist)

}

//int uint 根据操作系统位数   相当于int64  unint64
//字符byte占1位，  中文rune占3位
func TestIota(t *testing.T) {
	const (
		aaa = iota //默认0
		bbb        //每增加一行加一
	)
	fmt.Println(aaa, bbb)

	const (
		ccc = iota //默认0
		ddd = 100  //每增加一行加一
		fff = iota
	)
	fmt.Println(ccc, ddd, fff)
	int1 := 22
	fmt.Printf("%T \n", int1) //也可以打印类型
	println(&int1)
}
func TestMap(t *testing.T) { //map传递是一个东西
	map1 := make(map[string]string)
	map1["aa"] = "aa"
	dealmap(map1)
	fmt.Println(map1)
	slice := make([]string, 0, 3)
	slice = append(slice, "值1")
	dealslice(slice)
	fmt.Println(slice) //slice传递是一个东西
	chan1 := make(chan int, 10)
	chan1 <- 11
	chan1 <- 22
	dealChan(chan1) //chan 传递是一个东西
	for {
		select {
		case aa := <-chan1:
			fmt.Println("接一个", aa)
		default:
			fmt.Println("默认")
		}
	}

}
func dealmap(map1 map[string]string) {
	fmt.Println(map1)
	map1["aa"] = "bb"
}

func dealslice(slice1 []string) {
	fmt.Println(slice1)
	slice1[0] = "值2"
}

func dealChan(chan1 chan int) {
	chan1 <- 33
}

func TestPProf(t *testing.T) {
	file, _ := os.Create("E:\\tmp\\cpu.pprof")
	pprof.StartCPUProfile(file)
	defer pprof.StopCPUProfile()
	for i := 0; i < 10000; i++ {
		fmt.Println(i)
	}
}

//context 主线程 通知gorutine 停止
func TestContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go subGorutine(ctx)
	time.Sleep(5 * time.Second) //5秒去停gorutin
	cancel()
	chan1 := make(chan bool, 1)
	<-chan1
}
func subGorutine(ctx context.Context) {
ForLoop:
	for {
		fmt.Println("需要退出")
		time.Sleep(500 * time.Millisecond)
		select {
		case <-ctx.Done():
			fmt.Println("可以结束")
			break ForLoop
		default:

		}

	}
}
