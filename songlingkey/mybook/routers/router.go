package routers

import (
	"github.com/gin-gonic/gin"
	"songlingkey/mybook/pkg/setting"
	"songlingkey/mybook/routers/api/bookapi"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)

	apiGroup := r.Group("/api")
	{
		//获取书籍列表
		apiGroup.GET("/books", bookapi.GetBooks)

		apiGroup.POST("/addUser", bookapi.AddUser)

		//apiGroup.GET("/getUserByCode",bookapi.GetUserByCode)

		apiGroup.POST("/getUser", bookapi.GetUser)

		apiGroup.GET("/getStatusByUserId", bookapi.GetStatusByUserId)

		apiGroup.POST("/saveStudyHistory", bookapi.SaveStudyHistory)

		apiGroup.GET("/getContent", bookapi.GetContent)
	}

	return r
}
