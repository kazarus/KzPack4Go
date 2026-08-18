package KzPack4Go

import "fmt"

//import "github.com/dgrijalva/jwt-go"
import "github.com/golang-jwt/jwt/v5"

type TKzCustomClaims struct {
	WhoBuild int64   `json:"whoBuild"` //#租户实例
	UserRoot int64   `json:"userRoot"` //#租户实例
	EnrollID string  `json:"enrollID"` //#登录编码
	UserName string  `json:"userName"` //#用户名称
	EntityID string  `json:"entityID"` //#用户标识
	UserIndx int64   `json:"userIndx"` //#用户实例
	DataFrom int64   `json:"dataFrom"` //#请求终端
	CreateAt float64 `json:"createAt"` //#生成日期
	SynzFrom string  `json:"synzFrom"` //#同步方向
	jwt.RegisteredClaims
}

type TTokenUsr struct {
	WhoBuild int64   `json:"whoBuild"` //#租户实例
	UserRoot int64   `json:"userRoot"` //#租户实例
	EnrollID string  `json:"enrollID"` //#登录编码
	UserName string  `json:"userName"` //#用户名称
	EntityID string  `json:"entityID"` //#用户标识
	UserIndx int64   `json:"userIndx"` //#用户实例
	DataFrom int64   `json:"dataFrom"` //#请求终端
	CreateAt float64 `json:"createAt"` //#生成日期
	SynzFrom string  `json:"synzFrom"` //#同步方向

	Url4Path string `json:"url4Path"` //#请求链接
}

func ParseToken1(tokenString string, tokenSecret string) (*TKzCustomClaims, error) {

	KzSecret := []byte(tokenSecret)

	var KzClaims = TKzCustomClaims{}
	MyClaims, eror := jwt.ParseWithClaims(tokenString, &KzClaims, func(token *jwt.Token) (interface{}, error) {
		//#固定使用HMAC签名,防止算法混淆攻击
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return KzSecret, nil
	})

	if eror != nil {
		return nil, eror
	}

	if MyClaims != nil {
		if KzClaims, ok := MyClaims.Claims.(*TKzCustomClaims); ok && MyClaims.Valid {
			return KzClaims, nil
		}
	}

	return nil, jwt.ErrTokenInvalidClaims
}

func ParseToken2(tokenString string, tokenSecret string) (TTokenUsr, error) {

	KzSecret := []byte(tokenSecret)

	var KzClaims = TKzCustomClaims{}
	var TokenUsr = TTokenUsr{}
	MyClaims, eror := jwt.ParseWithClaims(tokenString, &KzClaims, func(token *jwt.Token) (interface{}, error) {
		//#固定使用HMAC签名,防止算法混淆攻击
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return KzSecret, nil
	})

	if eror != nil {
		return TokenUsr, eror
	}

	if MyClaims != nil {
		if KzClaims, ok := MyClaims.Claims.(*TKzCustomClaims); ok && MyClaims.Valid {

			TokenUsr.WhoBuild = KzClaims.WhoBuild
			TokenUsr.UserRoot = KzClaims.UserRoot
			TokenUsr.EnrollID = KzClaims.EnrollID
			TokenUsr.EntityID = KzClaims.EntityID
			TokenUsr.CreateAt = KzClaims.CreateAt
			TokenUsr.UserName = KzClaims.UserName
			TokenUsr.UserIndx = KzClaims.UserIndx
			TokenUsr.DataFrom = KzClaims.DataFrom
			TokenUsr.SynzFrom = KzClaims.SynzFrom

			return TokenUsr, nil
		}
	}

	return TokenUsr, jwt.ErrTokenInvalidClaims
}
