package handlers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/IamStubborN/filmtracker/database"

	"github.com/gin-gonic/gin"
)

func RespondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"error": message})
}

func RespondWithSuccess(c *gin.Context, code int, message interface{}) {
	c.JSON(code, gin.H{"success": message})
}

func userValidate(user *database.User) error {
	validateRegexp, err := regexp.Compile(`^[a-zA-Z]+(?:[_-]?[a-zA-Z0-9])*$`)
	if err != nil {
		log.Println(err)
	}
	if validateRegexp.MatchString(user.Login) &&
		validateRegexp.MatchString(user.Password) {
		return nil
	}
	return fmt.Errorf("login or password not valid")
}

func getUserFromContext(c *gin.Context) (*database.User, error) {
	user := database.User{}
	if err := c.BindJSON(&user); err != nil {
		return nil, err
	}
	if err := userValidate(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func InsertTokensToCookie(c *gin.Context, accessToken, refreshToken string, accessTokenExp, refreshTokenExp time.Time) {
	cookie := &http.Cookie{
		Name:     "Token",
		Value:    accessToken,
		Path:     "/",
		Expires:  accessTokenExp,
		HttpOnly: false,
	}
	http.SetCookie(c.Writer, cookie)
	cookie = &http.Cookie{
		Name:    "Refresh",
		Value:   refreshToken,
		Path:    "/",
		Expires: refreshTokenExp,

		HttpOnly: false,
	}
	http.SetCookie(c.Writer, cookie)
	//c.SetCookie(
	//	"Token",
	//	accessToken,
	//	int(time.Unix(0, int64(accessToken)).Second()),
	//	"/",
	//	"",
	//	false,
	//	false,
	//)
	//c.SetCookie(
	//	"Refresh",
	//	refreshToken,
	//	int(time.Duration(refreshTokenExp).Seconds()),
	//	"/",
	//	"",
	//	false,
	//	false,
	//)
}
func ClearCookie(c *gin.Context) {
	cookie := &http.Cookie{
		Name:    "Token",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),

		HttpOnly: false,
	}
	http.SetCookie(c.Writer, cookie)
	cookie = &http.Cookie{
		Name:    "Refresh",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),

		HttpOnly: false,
	}
	http.SetCookie(c.Writer, cookie)
	//c.SetCookie(
	//	"Token",
	//	"",
	//	-1,
	//	"/",
	//	"",
	//	false,
	//	false,
	//)
	//c.SetCookie(
	//	"Refresh",
	//	"",
	//	-1,
	//	"/",
	//	"",
	//	false,
	//	false,
	//)
}
