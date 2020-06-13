package bookapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-log/log"
	"github.com/medivhzhan/weapp"
	"net/http"
	"songlingkey/mybook/models"
)

func AddUser(c *gin.Context) {
	var basicUser = &models.BasicUser{} // 让所有人都可以听
	c.BindJSON(basicUser)
	fmt.Println("json得到", basicUser)
	openId, err := getOpenId(basicUser.UserCode)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "获取openid失败 ",
			"result":  nil,
		})
		return
	}
	basicUser.OpenId = openId
	models.AddUser(basicUser)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"msg":     "成功",
		"result":  basicUser,
	})
}

func GetStatusByUserId(c *gin.Context) {
	var returnStatus = http.StatusOK
	var userId = c.Query("userId")
	user, err := models.GetBasicUserByUserId(userId)
	if err != nil {
		returnStatus = http.StatusBadRequest
		user = &models.BasicUser{}
	}
	c.JSON(returnStatus, gin.H{
		"success": true,
		"msg":     "获取状态成功",
		"result":  user.Status,
	})
}

func GetUser(c *gin.Context) {
	var user = &models.BasicUser{}
	c.BindJSON(user)
	openId, err := getOpenId(user.UserCode)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "获取openid失败",
			"result":  nil,
		})
		return
	}
	user.OpenId = openId
	dbUser, recordNot := models.GetUserByOpenId(openId)
	//if(err!=nil){
	//	c.JSON(http.StatusOK, gin.H{
	//		"success" : false,
	//		"msg" : "根据openid获取失败",
	//		"result" : nil,
	//	})
	//	return
	//}
	if recordNot == true {
		models.AddUser(user)
	} else {
		user = dbUser
	}
	log.Log("返回的数据", user)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"msg":     "成功",
		"result":  user,
	})
}

func getOpenId(code string) (string, error) {
	res, err := weapp.Login("wx0722f8cf877f7db2", "352fa9156eba19d882470c51e16d7cac", code)
	if err != nil {
		// 处理一般错误信息
		return "", nil
	}
	fmt.Printf("返回结果: %#v", res.OpenID)

	return res.OpenID, nil

}

//func TestWXBizDataCrypt_Decrypt(entendUser *models.ExtendUser) {
//	appID := "wx0722f8cf877f7db2"
//	sessionKey := "tiihtNczf5v6AKRyjwEUhQ=="
//	encryptedData :=entendUser.EncryptedData
//	iv := entendUser.Iv
//	pc := util.NewWXBizDataCrypt(appID, sessionKey)
//	userInfo, err := pc.Decrypt(encryptedData, iv)
//	if err != nil {
//		fmt.Println(userInfo)
//
//	} else {
//		fmt.Println(userInfo)
//	}
//}
//func GetUserByCode(c *gin.Context) {
//	var code =c.Query("userCode")
//	book,_:=models.GetUserByCode(code)
//	c.JSON(http.StatusOK, gin.H{
//		"code" : e.SUCCESS,
//		"msg" : e.GetMsg(e.SUCCESS),
//		"result" : book,
//	})
//}
