package router

import (
	"github.com/IamStubborN/filmtracker/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CreateRouter() *gin.Engine {
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		middlewares.XSSMiddle(),
		cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost"},
			AllowHeaders:     []string{"Accept", "Content-Type"},
			AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
			AllowCredentials: true,
		}))
	AddV1RouterGroup(router)
	AddAuthRouterGroup(router)
	return router
}
