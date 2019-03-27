package handlers

import (
	"net/http"

	"github.com/IamStubborN/filmtracker/database"
	"github.com/IamStubborN/filmtracker/jwtmanager"

	"github.com/gin-gonic/gin"
)

var jmg = jwtmanager.GetJWTManager()
var db = database.GetDB()

func SignOut(c *gin.Context) {
	accessToken, err := c.Request.Cookie("Token")
	if err != nil {
		RespondWithError(c, http.StatusAccepted, err)
		return
	}
	refreshToken, err := c.Request.Cookie("Refresh")
	if err != nil {
		RespondWithError(c, http.StatusAccepted, err)
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

func SignIn(c *gin.Context) {
	userFromCtx, err := getUserFromContext(c)
	if err != nil {
		RespondWithError(c, http.StatusNotAcceptable, err.Error())
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
	} else {
		RespondWithError(c, http.StatusNotAcceptable, "this user isn't in the database.")
	}
}

func SignUp(c *gin.Context) {
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
