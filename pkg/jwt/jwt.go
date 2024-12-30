package jwt

import (
	"errors"
	"time"

	"github.com/spf13/viper"

	"go.uber.org/zap"

	"github.com/dgrijalva/jwt-go"
)

//const TokenExpireDuration = time.Hour * 2

var mySecret = []byte("漏船载酒")

// MyClaims 在官方基础上自定义字段
type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenToken 生成token
func GenToken(userID int64, username string) (string, error) {
	// 实例化一个声明
	c := MyClaims{
		userID,
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(viper.GetInt("auth.jwt_expire")) * time.Hour).Unix(),
			Issuer:    "Dandelion",
		},
	}
	// 使用声明和签名方法来生成签名对象token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定密钥mySecret获取编码后的字符串token
	signedString, err := token.SignedString(mySecret)
	if err != nil {
		zap.L().Error("GenToken/SignedString error: ", zap.Error(err))
	}
	return signedString, err
}

// 解析token
func ParseToken(tokenString string) (*MyClaims, error) {

	var myClaim = new(MyClaims)
	// Keyfunc函数，可以根据 Token 中的某个字段选择密钥；此处使用默认的匿名函数即可；
	// 因为就一个密钥，直接返回即可；如果有多个密钥，可以判断token中是否有指定字符串来选择
	token, err := jwt.ParseWithClaims(tokenString, myClaim, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return myClaim, nil
	}
	zap.L().Error("jwt parse token failed:", zap.Error(errors.New("invalid token")))
	return nil, errors.New("invalid token")
}
