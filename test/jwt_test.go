package test

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"testing"
	"time"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	jwt.StandardClaims
}

// TODO 签名密钥
var myKey = []byte("gin-gorm-oj-key")

func TestGenerateToken(t *testing.T) {
	UserClaims := &UserClaims{
		Identity: "user_1",
		Name:     "get",
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 60, // 令牌在当前时间的前60秒之前不生效
			ExpiresAt: time.Now().Unix() + 5,  // 令牌将在当前时间的后5秒过期
			Issuer:    "lsq",                  // 令牌的发行者
		},
	}
	//TODO  使用指定的签名方法和声明创建一个新的令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims)
	//TODO  使用签名密钥对令牌进行签名，并获取完整的签名后的令牌字符串
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(tokenString)
	//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6InVzZXJfMSIsIm5hbWUiOiJnZXQifQ.taKAmMZu6-6ioE4hSKqUgn9lHrqXSw-2TyEQTeNOreA
}

func TestAnalyseToken(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6InVzZXJfMSIsIm5hbWUiOiJnZXQifQ.taKAmMZu6-6ioE4hSKqUgn9lHrqXSw-2TyEQTeNOreA"
	userClaim := new(UserClaims)
	//TODO 解析 JWT，并在解析的过程中验证签名
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if claims.Valid {
		fmt.Println(userClaim)
	}
}
