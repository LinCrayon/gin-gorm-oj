package helper

import (
	"crypto/md5"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"time"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	IsAdmin  int    `json:"is_admin"`
	jwt.StandardClaims
}

func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// TODO 签名密钥
var myKey = []byte("gin-gorm-oj-key")

// GenerateToken
// TODO 生成token
func GenerateToken(identity, name string, isAdmin int) (string, error) {
	UserClaims := &UserClaims{
		Identity: "user_1",
		Name:     "get",
		IsAdmin:  isAdmin,
		StandardClaims: jwt.StandardClaims{
			//NotBefore: time.Now().Unix() - 60, // 令牌在当前时间的前60秒之前不生效
			//ExpiresAt: time.Now().Unix() + 5,  // 令牌将在当前时间的后5秒过期
			Issuer: "lsq", // 令牌的发行者
		},
	}
	//TODO  使用指定的签名方法和声明创建一个新的令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims)
	//TODO  使用签名密钥对令牌进行签名，并获取完整的签名后的令牌字符串
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}
	fmt.Println(tokenString)
	return tokenString, nil
	//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6InVzZXJfMSIsIm5hbWUiOiJnZXQifQ.taKAmMZu6-6ioE4hSKqUgn9lHrqXSw-2TyEQTeNOreA
}

// AnalyseToken
// 解析 token
func AnalyseToken(tokenString string) (*UserClaims, error) {
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("analyse Token Error:%v", err)
	}
	fmt.Println(claims)
	return userClaim, nil
}

// SendCode
// 发送验证码
func SendCode(toUserEmail, code string) error {
	e := email.NewEmail()
	e.From = "Crayon <2993373191@qq.com>"
	e.To = []string{toUserEmail}
	e.Subject = "验证码已发送，请查收"

	e.HTML = []byte("<h1>验证码:</h1>" + code)
	return e.Send("smtp.qq.com:587", smtp.PlainAuth("", "2993373191@qq.com", "vlkrwkqjayqedehc", "smtp.qq.com"))
}

// GetUUID
func GetUUID() string {
	return uuid.NewV4().String()
}

// 生成验证码
func GetRand() string {
	rand.Seed(time.Now().UnixNano())
	s := ""
	for i := 0; i < 6; i++ {
		s += strconv.Itoa(rand.Intn(10))
	}
	return s
}

// CodeSave
// 保存代码
func CodeSave(code []byte) (string, error) {
	dirName := "code/" + GetUUID()
	path := dirName + "/main.go"
	err := os.Mkdir(dirName, 0777)
	if err != nil {
		return "", err
	}
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	f.Write(code)
	defer f.Close()
	return path, nil
}
