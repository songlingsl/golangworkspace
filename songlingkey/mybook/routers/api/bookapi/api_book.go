package bookapi

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"songlingkey/mybook/models"
)

// 获取所有书籍
func GetBooks(c *gin.Context) {

	list := models.GetBooks()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"msg":     "获取所有书籍成功",
		//"msg":    "临时",
		"result": list,
	})
}

func GetContent(c *gin.Context) {
	var bookId = c.Query("bookId")
	result := models.GetContent(bookId)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"msg":     "获取内容成功",
		"result":  result,
	})
}
