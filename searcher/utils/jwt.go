package utils

import (
	"goSearcher/web/model"
	"time"

	"github.com/dgrijalva/jwt-go"
)

/**
Token工具类
token由三部分组成，
第一部分 协议头header 存储的是token使用的加密协议。
第二部分 token的内容，
第三部分 前面两部分+key一起哈希的值
*/

var jwtKey = []byte("a_secret_crect")

type Claims struct {
	UserId             uint
	jwt.StandardClaims //token内容描述信息
}

/**
根据user的id生成token
返回的是加密且签名后的token
**/
func ReleaseToken(user model.User) (string, error) {
	expirationTime := time.Now().Add(2 * time.Hour) //设置过期时间
	claims := &Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(), //token发放的时间
			Issuer:    "zjz",             //谁发放的token
			Subject:   "user token",      //主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //生成token
	tokenString, err := token.SignedString(jwtKey)             //对token签名

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

/**
解析token 返回token信息，和对token的补充信息（claims）
**/
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {

		return jwtKey, nil
	})

	return token, claims, err
}
