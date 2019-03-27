package router

import (
	"github.com/IamStubborN/filmtracker/handlers"
	"github.com/gin-gonic/gin"
)

func AddAuthRouterGroup(router *gin.Engine) {
	auth := router.Group("/users/auth/")
	auth.Use()
	auth.POST("signup/", handlers.SignUp)
	auth.POST("signin/", handlers.SignIn)
	auth.GET("signout/", handlers.SignOut)
}
