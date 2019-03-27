package router

import (
	"github.com/IamStubborN/filmtracker/handlers"
	"github.com/IamStubborN/filmtracker/middlewares"
	"github.com/gin-gonic/gin"
)

func AddV1RouterGroup(router *gin.Engine) {
	v1 := router.Group("/api/v1/")
	v1.Use(
		middlewares.JWTMiddleware(),
	)
	v1.GET("/", handlers.StartPage)
	v1.GET("films/", handlers.FetchAllFilms)
	v1.GET("films/filter", handlers.FetchFilmsFilter)
	v1.POST("films/", middlewares.AccessMiddleware(), handlers.AddFilm)
	v1.PUT("films/", middlewares.AccessMiddleware(), handlers.UpdateFilm)
	v1.GET("films/film/:id", handlers.FetchSingleFilm)
	v1.DELETE("films/film/:id", middlewares.AccessMiddleware(), handlers.DeleteFilm)
	v1.GET("films/page/:id", handlers.FetchPage)
}
