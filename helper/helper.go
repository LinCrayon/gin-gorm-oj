package helper

import (
	"crypto/md5"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"sync"
	"time"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	IsAdmin  int    `json:"is_admin"`
	jwt.StandardClaims
}

// Snowflake算法
var (
	// 设置Snowflake算法中的起始时间点,这个时间戳表示了当前时间与指定的起始时间之间的毫秒数
	epoch = time.Date(2024, time.January, 01, 00, 00, 00, 00, time.UTC).UnixMilli()
)

const (
	timestampBits    = 41 // 时间戳占用的位数
	dataCenterIdBits = 5  // 数据中心ID占用的位数
	workerIdBits     = 5  // 工作ID（或机器ID）占用的位数
	seqBits          = 12 //序列号占用的位数，用于标识同一毫秒内生成的不同ID的序列

	// 时间戳的最大值, just like 2^41-1 = 2199023255551
	timestampMaxValue = -1 ^ (-1 << timestampBits)
	// dataCenterId max value, just like 2^5-1 = 31
	dataCenterIdMaxValue = -1 ^ (-1 << dataCenterIdBits)
	// workId max value, just like 2^5-1 = 31
	workerIdMaxValue = -1 ^ (-1 << workerIdBits)
	// 序列号的最大值, just like 2^12-1 = 4095
	seqMaxValue = -1 ^ (-1 << seqBits)

	workIdShift       = 12 // number of workId offsets (seqBits)
	dataCenterIdShift = 17 // number of dataCenterId offsets (seqBits + workerIdBits)
	timestampShift    = 22 // number of timestamp offsets (seqBits + workerIdBits + dataCenterIdBits)

	defaultInitValue = 0
)

type SnowflakeSeqGenerator struct {
	mu           *sync.Mutex
	timestamp    int64
	dataCenterId int64
	workerId     int64
	sequence     int64
}

// NewSnowflakeSeqGenerator Snowflake算法
func NewSnowflakeSeqGenerator(dataCenterId, workId int64) (r *SnowflakeSeqGenerator, err error) {
	if dataCenterId < 0 || dataCenterId > dataCenterIdMaxValue {
		err = fmt.Errorf("dataCenterId should between 0 and %d", dataCenterIdMaxValue-1)
		return
	}
	if workId < 0 || workId > workerIdMaxValue {
		err = fmt.Errorf("workId should between 0 and %d", dataCenterIdMaxValue-1)
		return
	}
	return &SnowflakeSeqGenerator{
		mu:           new(sync.Mutex),      //创建一个新的互斥锁
		timestamp:    defaultInitValue - 1, //默认初始值减去1，以确保在第一次生成ID时可以更新时间戳。
		dataCenterId: dataCenterId,
		workerId:     workId,
		sequence:     defaultInitValue, //同一毫秒内的序列号
	}, nil
}
func (S *SnowflakeSeqGenerator) GenerateId(entity string, ruleName string) string {
	S.mu.Lock()
	defer S.mu.Unlock()

	now := time.Now().UnixMilli()
	if S.timestamp > now { // Clock callback 时钟发生回拨
		log.Printf("Clock moved backwards. Refusing to generate ID, last timestamp is %d, now is %d", S.timestamp, now)
		return ""
	}
	if S.timestamp == now { // generate multiple IDs in the same millisecond, incrementing the sequence number to prevent conflicts
		S.sequence = (S.sequence + 1) & seqMaxValue //位与,可以确保序列号始终保持在合法的范围内
		if S.sequence == 0 {                        //如果序列号已经达到最大值（等于0），则进入循环等待下一毫秒
			for now <= S.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else { // initialized sequences are used directly at different millisecond timestamps
		S.sequence = defaultInitValue
	}
	tmp := now - epoch
	if tmp > timestampMaxValue {
		log.Printf("epoch should between 0 and %d", timestampMaxValue-1)
		return ""
	}
	S.timestamp = now
	// combine the parts to generate the final ID and convert the 64-bit binary to decimal digits.
	r := (tmp)<<timestampShift | //将时间戳部分左移 timestampShift 位
		(S.dataCenterId << dataCenterIdShift) | //将数据中心ID左移 dataCenterIdShift 位
		(S.workerId << workIdShift) | //将工作ID左移 workIdShift 位
		(S.sequence) //直接使用序列号部分
	return fmt.Sprintf("%d", r)
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
	return tokenString, nil
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
	//smtp.PlainAuth身份验证
	return e.Send("smtp.qq.com:587", smtp.PlainAuth("", "2993373191@qq.com", "vlkrwkqjayqedehc", "smtp.qq.com"))
}

func GetUUID() string {
	return uuid.NewV4().String()
}

// GetRand 生成验证码
func GetRand() string {
	rand.Seed(time.Now().UnixNano()) //随机数生成器的种子
	s := ""
	for i := 0; i < 6; i++ {
		s += strconv.Itoa(rand.Intn(10)) //rand.Intn(10) 来生成一个介于 0 和 9 之间的随机整数
	}
	return s
}

// CodeSave
// 保存代码
func CodeSave(code []byte) (string, error) {
	var dataCenterId, workId int64 = 1, 1
	generator, err := NewSnowflakeSeqGenerator(dataCenterId, workId)
	if err != nil {
		log.Fatal("NewSnowflakeSeqGenerator Error")
	}
	dirName := "code/" + generator.GenerateId("", "")
	path := dirName + "/main.go"
	err = os.Mkdir(dirName, 0777)
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
