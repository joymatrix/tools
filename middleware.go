package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func securityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//判断token是否过期，过期就重新登录
		userId, exists := c.Get("userId")
		if !exists {
			c.Abort()
			failWithMsg(c, constants.ErrToken.Code, constants.ErrToken.Msg)
			return
		}
		token, exists := c.Get("token")
		if !exists {
			c.Abort()
			failWithMsg(c, constants.ErrToken.Code, constants.ErrToken.Msg)
			return
		}
		strUserId := strconv.FormatInt(userId.(int64), 10)
		oldToken, err := utils.GetToken(c, strUserId)
		if err != nil {
			utils.GetLog().Errorf("token is expired err:%+v", err.Error())
			c.Abort()
			failWithMsg(c, constants.ErrToken.Code, constants.ErrToken.Msg)
			return
		}

		if token != oldToken {
			utils.GetLog().Error("token is invalid")
			c.Abort()
			failWithMsg(c, constants.ErrToken.Code, constants.ErrToken.Msg)
			return
		}
		expireTime := 24 * 3600
		strUserId = strconv.FormatInt(userId.(int64), 10)
		utils.SetToken(c, strUserId, token.(string), int64(expireTime))
		c.SetCookie("token", token.(string), expireTime, "/sanicalc/", "", false, false)
		c.Next()
	}
}

func prepareMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//解析token 设置上下文
		token := c.GetHeader("Authentication")
		if token == "" {
			utils.GetLog().Info("token empty")
			c.Abort()
			failWithMsg(c, constants.ErrToken.Code, constants.ErrToken.Msg)
			return
		}

		err := utils.ParseToken(c, token)
		if err != nil {
			utils.GetLog().Infof("parse token err:%+v", err.Error())
			c.Abort()
			failWithMsg(c, constants.ErrToken.Code, constants.ErrToken.Msg)
			return
		}

		c.Next()
	}
}
