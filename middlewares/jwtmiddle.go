package middlewares

import (
	"net/http"

	"github.com/IamStubborN/filmtracker/handlers"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//func JWTMiddleware1() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		accessTokenCookie, err := c.Cookie("Token")
//		if err != nil {
//			accessTokenCookie = ""
//		}
//
//	}
//}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessTokenCookie, err := c.Cookie("Token")
		if err != nil {
			accessTokenCookie = ""
		}
		refreshTokenCookie, err := c.Cookie("Refresh")
		if err != nil {
			handlers.RespondWithError(c, http.StatusNotAcceptable, "no Refresh cookie")
			return
		}
		if jmg.IsTokenInBlackList(accessTokenCookie) {
			handlers.RespondWithError(c, http.StatusNotAcceptable, "invalid token, please login or sign up")
			return
		}
		if accessTokenCookie != "" {
			accessToken, err := jmg.ParseToken(accessTokenCookie)
			if err != nil {
				refreshTokenLogic(c, refreshTokenCookie)
			}
			if accessToken.Valid {
				c.Next()
			}
		} else if accessTokenCookie == "" {
			refreshTokenLogic(c, refreshTokenCookie)
		} else {
			handlers.RespondWithError(c, http.StatusUnauthorized, "API token required")
		}
	}
}

func refreshTokenLogic(c *gin.Context, refreshTokenCookie string) {
	refreshToken, err := jmg.ParseToken(refreshTokenCookie)
	if err != nil {
		handlers.RespondWithError(c, http.StatusUnauthorized, err.Error())
		return
	}
	if !refreshToken.Valid {
		handlers.RespondWithError(c, http.StatusUnauthorized, "refresh Token invalid")
		return
	}
	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		handlers.RespondWithError(c, http.StatusInternalServerError, "can't parse accessToken claims")
		return
	}
	UserID := claims["user_id"].(string)
	user, err := db.GetUserByUserID(UserID)
	if err != nil {
		handlers.RespondWithError(c, http.StatusUnauthorized, "User isn't exist")
		return
	}
	if !(refreshToken.Raw == user.RefreshToken) {
		handlers.RespondWithError(c, http.StatusUnauthorized, "Bad refresh Token")
		return
	}
	newAccessTokenExp := jmg.GenExp("access")
	newRefreshTokenExp := jmg.GenExp("refresh")
	newRefreshToken, err := jmg.GenerateToken(user.UserID, user.Role, newRefreshTokenExp)
	if err != nil {
		handlers.RespondWithError(c, http.StatusUnauthorized, err)
		return
	}
	newAccessToken, err := jmg.GenerateToken(user.UserID, user.Role, newAccessTokenExp)
	if err != nil {
		handlers.RespondWithError(c, http.StatusUnauthorized, err)
		return
	}
	if err := db.UpdateRefreshTokenForUser(user.UserID, newRefreshToken); err != nil {
		handlers.RespondWithError(c, http.StatusUnauthorized, err)
		return
	}
	handlers.InsertTokensToCookie(c, newAccessToken, newRefreshToken, newAccessTokenExp, newRefreshTokenExp)
}
