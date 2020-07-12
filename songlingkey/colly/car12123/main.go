package main

import (
	"github.com/gocolly/colly"
	"github.com/muyids/chrome-cookie"
	"log"
	"time"
)

func main() {
	chrome_cookie.GetCookie("https://cloud.baidu.com/", "C:\\Users\\songlinga\\AppData\\Local\\Google\\Chrome\\User Data\\Default\\Cookies", "json")
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36"),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Content-Type", "application/x-www-form-urlencoded")
		r.Headers.Set("Cookie", "_uab_collina=159436047589222596954183; TS011422ee=015e99a05df7de552bf53c7b6f0b4016988fd4dfb364542b3e17fca98915cb188542d213cd701b591c535faf10de2cdad1db0672ad; insert_cookie=67313298; JSESSIONID-L=a5522ecf-114b-4fe5-ac61-1cec4d25808a; accessToken=izseEXD3umt53giVGUK1wJmOcIdK+hXGfu/xuiCb+iX33uiRoV3dH9tCqdKlHrrmb2pEjERfCJun/iR9I0bssjoQFgO+7dezE+MMo1hxeNctvma2/TjrJT+g5gnrIp0Tf9ac1jjrDkYE8s/0KWkxSgFUDlAzAGLS0tbjaXIAHEW9OkEx1gYt+dMH03BxbXnO; 122_TAG_SHOWDWINDEXPUB=1")

	})

	c.OnResponse(func(r *colly.Response) {

		log.Println("内容是：", string(r.Body))

	})

	for {
		println("循环读取。。。")
		err := c.Post("https://bj.122.gov.cn/veh/jcz/yyxx/list", map[string]string{"ksrq": "2020-07-10", "jsrq": "2020-07-17"})
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(5 * time.Minute)
	}

	c.Wait()

}
