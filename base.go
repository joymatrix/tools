package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func respWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": constants.SuccessMsg.Code,
		"msg":  "",
		"data": data,
	})
}

func success(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": constants.SuccessMsg.Code,
		"msg":  "",
	})
}

func failWithMsg(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code": code,
		"msg":  msg,
	})
}
