package router

import (
	"github.com/IamStubborN/filmtracker/middlewares"
	"github.com/gin-gonic/gin"
)

func CreateRouter() *gin.Engine {
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		middlewares.XSSMiddle(),
	)
	AddV1RouterGroup(router)
	AddAuthRouterGroup(router)
	return router
}
