package test

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	_ "github.com/ugorji/go/codec"
	"testing"
)

var (
	ErrAppIDNotMatch       = errors.New("app id not match")
	ErrInvalidBlockSize    = errors.New("invalid block size")
	ErrInvalidPKCS7Data    = errors.New("invalid PKCS7 data")
	ErrInvalidPKCS7Padding = errors.New("invalid padding on input")
)

type UserInfo struct {
	OpenID    string `json:"openId"`
	UnionID   string `json:"unionId"`
	NickName  string `json:"nickName"`
	Gender    int    `json:"gender"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	AvatarURL string `json:"avatarUrl"`
	Language  string `json:"language"`
	Watermark struct {
		Timestamp int64  `json:"timestamp"`
		AppID     string `json:"appid"`
	} `json:"watermark"`
}

type WXBizDataCrypt struct {
	appID, sessionKey string
}

func NewWXBizDataCrypt(appID, sessionKey string) *WXBizDataCrypt {
	return &WXBizDataCrypt{
		appID:      appID,
		sessionKey: sessionKey,
	}
}

// pkcs7Unpad returns slice of the original data without padding
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if len(data)%blockSize != 0 || len(data) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	c := data[len(data)-1]
	n := int(c)
	if n == 0 || n > len(data) {
		return nil, ErrInvalidPKCS7Padding
	}
	for i := 0; i < n; i++ {
		if data[len(data)-n+i] != c {
			return nil, ErrInvalidPKCS7Padding
		}
	}
	return data[:len(data)-n], nil
}

func (w *WXBizDataCrypt) Decrypt(encryptedData, iv string) (*UserInfo, error) {
	aesKey, err := base64.StdEncoding.DecodeString(w.sessionKey)
	if err != nil {
		return nil, err
	}
	cipherText, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}
	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, ivBytes)
	mode.CryptBlocks(cipherText, cipherText)
	cipherText, err = pkcs7Unpad(cipherText, block.BlockSize())
	if err != nil {
		return nil, err
	}
	var userInfo UserInfo
	err = json.Unmarshal(cipherText, &userInfo)
	if err != nil {
		return nil, err
	}
	if userInfo.Watermark.AppID != w.appID {
		return nil, ErrAppIDNotMatch
	}
	return &userInfo, nil
}





func TestWXBizDataCrypt_Decrypt_test(t *testing.T) {
	appID := "wx4f4bc4dec97d474b"
	sessionKey := "tiihtNczf5v6AKRyjwEUhQ=="
	encryptedData := "CiyLU1Aw2KjvrjMdj8YKliAjtP4gsMZMQmRzooG2xrDcvSnxIMXFufNstNGTyaGS9uT5geRa0W4oTOb1WT7fJlAC+oNPdbB+3hVbJSRgv+4lGOETKUQz6OYStslQ142dNCuabNPGBzlooOmB231qMM85d2/fV6ChevvXvQP8Hkue1poOFtnEtpyxVLW1zAo6/1Xx1COxFvrc2d7UL/lmHInNlxuacJXwu0fjpXfz/YqYzBIBzD6WUfTIF9GRHpOn/Hz7saL8xz+W//FRAUid1OksQaQx4CMs8LOddcQhULW4ucetDf96JcR3g0gfRK4PC7E/r7Z6xNrXd2UIeorGj5Ef7b1pJAYB6Y5anaHqZ9J6nKEBvB4DnNLIVWSgARns/8wR2SiRS7MNACwTyrGvt9ts8p12PKFdlqYTopNHR1Vf7XjfhQlVsAJdNiKdYmYVoKlaRv85IfVunYzO0IKXsyl7JCUjCpoG20f0a04COwfneQAGGwd5oa+T8yO5hzuyDb/XcxxmK01EpqOyuxINew=="
	iv := "r7BXXKkLb8qrSNn05n0qiA=="
	pc := NewWXBizDataCrypt(appID, sessionKey)
	userInfo, err := pc.Decrypt(encryptedData, iv)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%#v\n", userInfo)
	}
}