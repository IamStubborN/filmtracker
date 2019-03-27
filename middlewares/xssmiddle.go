package middlewares

import (
	"github.com/dvwright/xss-mw"
	"github.com/gin-gonic/gin"
)

func XSSMiddle() gin.HandlerFunc {
	xssRemover := xss.XssMw{}
	return xssRemover.RemoveXss()
}
