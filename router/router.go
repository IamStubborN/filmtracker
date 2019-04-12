package router

import (
	"github.com/IamStubborN/filmtracker/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func CreateRouter() *gin.Engine {
	router := gin.New()
	router.Use(
		gin.Recovery(),
		middlewares.XSSMiddle(),
		cors.New(cors.Config{
			AllowOrigins: []string{
				"http://filmtracker-api.com",
				"http://localhost",
				"http://localhost:3000",
				"http://localhost:5555",
			},
			AllowHeaders:     []string{"Accept", "Content-Type"},
			AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
			AllowCredentials: true,
		}))
	AddV1RouterGroup(router)
	AddAuthRouterGroup(router)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
