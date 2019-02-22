package impl

import (
	"crypto/rsa"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"tianwei.pro/business"
	"tianwei.pro/micro/conf"
	"tianwei.pro/sam/core/dto"
	"tianwei.pro/sam/core/facade"
	"time"
)

var (
	tokenFacadeImpl facade.TokenFacade
	TokenExpired    = errors.New("token已过期,请重新登录")
	TokenInvalid    = errors.New("token验证失败")
)

type SamClaims struct {
	jwt.StandardClaims

	UserName string

	// 给哪个端的
	Terminal int8
}

func init() {
	if b, err := ioutil.ReadFile("conf/sam_pri"); err != nil {
		panic(err)
	} else if priv, err := jwt.ParseRSAPrivateKeyFromPEM(b); err != nil {
		panic(err)
	} else if b, err = ioutil.ReadFile("conf/sam_pub"); err != nil {
		panic(err)
	} else if pub, err := jwt.ParseRSAPublicKeyFromPEM(b); err != nil {
		panic(err)
	} else {
		tokenFacadeImpl = &TokenFacadeImpl{
			priv: priv,
			pub:  pub,
		}
	}
}

type TokenFacadeImpl struct {
	priv *rsa.PrivateKey
	pub  *rsa.PublicKey
}

func (t *TokenFacadeImpl) DecodeToken(tokenString string) (*dto.UserDto, error) {
	if token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return t.pub, nil
	}); err != nil {
		if err.Error() == "Token is expired" {
			return nil, TokenExpired
		} else if err.Error() == "token is invalid" {
			return nil, TokenInvalid
		}
		return nil, TokenInvalid
	} else if tokenMap, ok := token.Claims.(jwt.MapClaims); !ok {
		logs.Error("parse token switch map claim failed. %v", token.Claims)
		return nil, TokenInvalid
	} else {
		id := business.CastStringToInt64(tokenMap["jti"].(string))
		if id == 0 {
			return nil, TokenInvalid
		}

		return &dto.UserDto{
			BaseDto: dto.BaseDto{
				Id: id,
			},
			UserName: tokenMap["UserName"].(string),
		}, nil
	}
}

func (t *TokenFacadeImpl) EncodeToken(user *dto.UserDto, Terminal int8) (string, error) {

	// 默认30天过期
	expiredTime := time.Duration(conf.Conf.DefaultInt64("tokenExpire", 24*60*30))
	s := jwt.StandardClaims{
		Id:        business.CastInt64ToString(user.Id),
		ExpiresAt: time.Now().Add(time.Minute * expiredTime).Unix(),
	}

	samClaims := &SamClaims{
		StandardClaims: s,
		UserName:       user.UserName,
		Terminal:       Terminal,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS384, samClaims)

	if ss, err := token.SignedString(t.priv); err != nil {
		logs.Error("signed token failed. user: %v, err: %v", user, err)
		return "", SystemError
	} else {
		return ss, nil
	}
}
