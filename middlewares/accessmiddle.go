package middlewares

import (
	"net/http"

	"github.com/IamStubborN/filmtracker/database"
	"github.com/IamStubborN/filmtracker/handlers"
	"github.com/IamStubborN/filmtracker/jwtmanager"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jmg *jwtmanager.JwtManager
var db *database.Database

func init() {
	jmg = jwtmanager.GetJWTManager()
	db = database.GetDB()
}

func AccessMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessTokenCookie, err := c.Cookie("Token")
		if err != nil {
			handlers.RespondWithError(c, http.StatusNotAcceptable, err.Error())
			return
		}
		refreshTokenCookie, err := c.Cookie("Refresh")
		if err != nil {
			handlers.RespondWithError(c, http.StatusNotAcceptable, err.Error())
			return
		}
		token, err := jmg.ParseToken(accessTokenCookie)
		if err != nil {
			handlers.RespondWithError(c, http.StatusNotAcceptable, err.Error())
			return
		}
		if token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				UserID := claims["user_id"].(string)
				user, err := db.GetUserByUserID(UserID)
				if err != nil {
					handlers.RespondWithError(c, http.StatusNotAcceptable, "User isn't exist")
				}
				if user.Role == "admin" && user.RefreshToken == refreshTokenCookie {
					c.Next()
				} else {
					handlers.RespondWithError(c, http.StatusUnauthorized, "Access denied")
				}
			}
		}
	}
}
