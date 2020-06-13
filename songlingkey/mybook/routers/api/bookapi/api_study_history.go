package bookapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"songlingkey/mybook/models"
)

func SaveStudyHistory(c *gin.Context) {
	var studyHistory = &models.StudyHistory{}
	c.BindJSON(studyHistory)

	affected := models.UpdateHistory(studyHistory)
	fmt.Println("更新了？", affected)
	if affected == 0 {
		models.SaveHistory(studyHistory)
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"msg":     "成功",
		"result":  nil,
	})
}
