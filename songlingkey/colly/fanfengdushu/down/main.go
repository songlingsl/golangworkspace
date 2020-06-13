package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"time"
)

func main() {
	c := colly.NewCollector(
		colly.MaxBodySize(0),

	)

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("aaa")
		r.Save("e://tmp//wo123.mp3")
	})

	c.SetRequestTimeout(20*time.Second)
	//c.Visit("https://cdn-ali.dushu.io/audio/full/d2c3b76927a65547e81cf140e0d1e625_4h7xd6.mp3")
   c.Visit("https://cdn-ali-dest.dushu.io/media/audio/1567183023bc7b0b8950feaec12a3e358eded2d075nppaa4.mp3")

	c.Wait()
	forever := make(chan bool)
	<-forever
}
