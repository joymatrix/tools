package utils

import (
	"errors"
	"sanicalc/internal/constants"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func GenToken(phoneNum, account string, userId, expireTime int64) (token string, err error) {
	claims := jwt.MapClaims{
		"account":  account,
		"phoneNum": phoneNum,
		"userId":   userId,
		"iss":      "sanicalc",
		"nbf":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Duration(expireTime) * time.Second).Unix(),
		"iat":      time.Now().Unix(),
	}

	claimsObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err = claimsObj.SignedString([]byte(constants.JWT_SIGN_KEY))

	return
}

func secret() jwt.Keyfunc {
	return func(t *jwt.Token) (interface{}, error) {
		return []byte(constants.JWT_SIGN_KEY), nil
	}
}

func ParseToken(ctx *gin.Context, token string) (err error) {
	tokn, err := jwt.Parse(token, secret())
	if err != nil {
		return
	}

	claim, ok := tokn.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("parse error")
		return
	}

	if !tokn.Valid {
		err = errors.New("token error")
		return
	}

	account, ok := claim["account"].(string)
	if !ok {
		err = errors.New("account error")
		return
	}
	phoneNum, ok := claim["phoneNum"].(string)
	if !ok {
		err = errors.New("phoneNum error")
		return
	}

	userId, ok := claim["userId"].(float64)
	if !ok {
		err = errors.New("userId error")
		return
	}

	ctx.Set("userId", int64(userId))
	ctx.Set("account", account)
	ctx.Set("phoneNum", phoneNum)
	ctx.Set("token", token)

	return
}
