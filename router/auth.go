package router

import (
	"github.com/IamStubborN/filmtracker/handlers"
	"github.com/gin-gonic/gin"
)

func AddAuthRouterGroup(router *gin.Engine) {
	auth := router.Group("/users/auth/")
	auth.Use()
	auth.OPTIONS("signup/", handlers.SignUp)
	auth.POST("signup/", handlers.SignUp)
	auth.OPTIONS("signin/", handlers.SignIn)
	auth.POST("signin/", handlers.SignIn)
	auth.GET("signout/", handlers.SignOut)
}
