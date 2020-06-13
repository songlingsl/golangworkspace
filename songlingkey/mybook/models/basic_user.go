package models

import "github.com/prometheus/common/log"

type BasicUser struct {
	UserId    int `gorm:"PRIMARY_KEY:user_id"`
	Gender    int
	NickName  string `json:"nickName"`
	Province  string
	City      string
	AvatarUrl string
	Status    int
	Phone     string
	UserCode  string `gorm:"-"`
	OpenId    string
}

//type ExtendUser struct{
//	UserInfo BasicUser
//	EncryptedData string
//	Iv string
//
//}

func AddUser(user *BasicUser) error {
	user.Status = 1
	log.Info("保存用户", user)
	if err := db.Create(user).Error; err != nil {
		return err
	}

	return nil
}

func GetUserByOpenId(openId string) (*BasicUser, bool) {
	var user BasicUser

	if recordNot := db.First(&user, "open_id = ?", openId).RecordNotFound(); recordNot == true {
		return nil, recordNot

	}

	return &user, false
}
func GetBasicUserByUserId(userId string) (*BasicUser, error) {
	var user BasicUser
	if err := db.First(&user, userId).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
