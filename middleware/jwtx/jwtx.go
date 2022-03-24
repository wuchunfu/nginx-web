package jwtx

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/golang-jwt/jwt"
	"strconv"
	"time"
)

var jwtKey = "jwt_secret_key"

type MyClaims struct {
	UserId   int64
	Username string
	jwt.StandardClaims
}

// GenerateToken 生成 jwt token
func GenerateToken(userId int64, username string) (map[string]string, error) {
	flag := Md5(userId)
	// 过期时间为七天
	expireTime := time.Now().Add(7 * 24 * time.Hour).Unix()
	claims := &MyClaims{
		UserId:   userId,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			Id:        flag,              // 唯一身份标识
			Subject:   "User token",      // 主题
			IssuedAt:  time.Now().Unix(), // 颁发时间
			Issuer:    username,          // 颁发者
			NotBefore: time.Now().Unix(), // 生效时间
			ExpiresAt: expireTime,        // 过期时间
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := tokenClaims.SignedString([]byte(jwtKey))
	if err != nil {
		return nil, err
	}

	tokenMap := make(map[string]string, 10)
	tokenMap["userId"] = strconv.FormatInt(userId, 10)
	tokenMap["token"] = signedToken
	tokenMap["expireTime"] = time.Unix(expireTime, 0).Format("2006-01-02 15:04:05")
	return tokenMap, nil
}

// ParseToken 解析Token
func ParseToken(signedToken string) (*MyClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(signedToken, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*MyClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

// Md5 根据id使用md5生成用户唯一标识
func Md5(userId int64) string {
	md5Hash := md5.New()
	md5Hash.Write([]byte(strconv.FormatInt(userId, 10)))
	return hex.EncodeToString(md5Hash.Sum([]byte("")))
}
