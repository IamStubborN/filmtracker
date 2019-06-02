package handlers

import (
	"net/http"

	"github.com/IamStubborN/filmtracker/database"
	"github.com/IamStubborN/filmtracker/jwtmanager"

	"github.com/gin-gonic/gin"
)

var jmg = jwtmanager.GetJWTManager()
var db = database.GetDB()

// SignOut godoc
// @Summary Sign Out
// @Description Sign Out from server, delete cookies
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 "{"success":"Sign Out"}"
// @Failure 406 "{"error":"http: named cookie not present"/"wrong refresh token"}"
// @Router /users/auth/signout/ [get]
func SignOut(c *gin.Context) {
	accessToken, err := c.Request.Cookie("Token")
	if err != nil {
		RespondWithError(c, http.StatusNotAcceptable, err.Error())
		return
	}
	refreshToken, err := c.Request.Cookie("Refresh")
	if err != nil {
		RespondWithError(c, http.StatusNotAcceptable, err.Error())
		return
	}
	jmg.AddTokenToBlackList(accessToken.Value)
	if err := db.DeleteRefreshToken(refreshToken.Value); err != nil {
		RespondWithError(c, http.StatusNotAcceptable, "wrong refresh token")
		return
	}
	ClearCookie(c)
	RespondWithSuccess(c, http.StatusOK, "Sign out")
}

// SignIn godoc
// @Summary Sign In
// @Description Sign In into server, add cookies
// @Tags users
// @Accept  json
// @Produce  json
// @Param Login body database.User false "Add login and password"
// @Success 200 "{"success":"Sign In"}"
// @Header 200 {string} Token "JWT Token"
// @Header 200 {string} Refresh "JWT refresh Token"
// @Failure 406 "{"error":"this user isn't in the database."}"
// @Router /users/auth/signin/ [post]
func SignIn(c *gin.Context) {
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	userFromCtx, err := getUserFromContext(c)
	if err != nil {
		RespondWithError(c, http.StatusNotAcceptable, err.Error())
		return
	}
	userExist, err := db.IsExistUser(userFromCtx.Login)
	if err != nil {
		RespondWithError(c, http.StatusNotAcceptable, err.Error())
		return
	}
	if userExist {
		user, err := db.GetUser(userFromCtx.Login, userFromCtx.Password)
		if err != nil {
			RespondWithError(c, http.StatusNotAcceptable, err.Error())
			return
		}
		accessExp := jmg.GenExp("access")
		refreshExp := jmg.GenExp("refresh")
		accessToken, err := jmg.GenerateToken(user.UserID, user.Role, accessExp)
		if err != nil {
			RespondWithError(c, http.StatusNotAcceptable, err.Error())
			return
		}
		refreshToken, err := jmg.GenerateToken(user.UserID, user.Role, refreshExp)
		if err != nil {
			RespondWithError(c, http.StatusNotAcceptable, err.Error())
			return
		}
		if err := db.UpdateRefreshTokenForUser(user.UserID, refreshToken); err != nil {
			RespondWithError(c, http.StatusNotAcceptable, err.Error())
			return
		}
		InsertTokensToCookie(c, accessToken, refreshToken, accessExp, refreshExp)
		RespondWithSuccess(c, http.StatusOK, "Sign in")
		return
	} else {
		RespondWithError(c, http.StatusNotAcceptable, "this user isn't in the database.")
		return
	}
}

// SignUp godoc
// @Summary Sign Up
// @Description Sign Up into server, add cookies
// @Tags users
// @Accept  json
// @Produce  json
// @Param Login body database.User false "Add login and password"
// @Header 200 {string} Token "JWT Token"
// @Header 200 {string} Refresh "JWT refresh Token"
// @Success 200 "{"success":"Sign Up"}"
// @Failure 406 "{"error":"this user is already exist in database."}"
// @Router /users/auth/signup/ [post]
func SignUp(c *gin.Context) {
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	}
	user, err := getUserFromContext(c)
	if err != nil {
		RespondWithError(c, http.StatusNotAcceptable, err.Error())
		return
	}
	userExist, err := db.IsExistUser(user.Login)
	if err != nil {
		RespondWithError(c, http.StatusNotAcceptable, err.Error())
		return
	}
	if !userExist {
		user, err := db.CreateUser(user.Login, user.Password)
		if err != nil {
			RespondWithError(c, http.StatusNotAcceptable, err.Error())
			return
		}
		accessExp := jmg.GenExp("access")
		refreshExp := jmg.GenExp("refresh")
		accessToken, err := jmg.GenerateToken(user.UserID, user.Role, accessExp)
		if err != nil {
			RespondWithError(c, http.StatusNotAcceptable, err.Error())
			return
		}
		refreshToken, err := jmg.GenerateToken(user.UserID, user.Role, refreshExp)
		if err != nil {
			RespondWithError(c, http.StatusNotAcceptable, err.Error())
			return
		}
		if err := db.UpdateRefreshTokenForUser(user.UserID, refreshToken); err != nil {
			RespondWithError(c, http.StatusNotAcceptable, err.Error())
			return
		}
		InsertTokensToCookie(c, accessToken, refreshToken, accessExp, refreshExp)
		RespondWithSuccess(c, http.StatusOK, "Sign up")
	} else {
		RespondWithError(c, http.StatusNotAcceptable, "this user is already exist in database.")
	}
}
