package jwtmanager

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"github.com/dgrijalva/jwt-go"
)

type (
	JwtManager struct {
		AccessTokenLife      time.Duration
		RefreshTokenLife     time.Duration
		accessTokenBlackList map[string]struct{}
	}
)

var jmg *JwtManager

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("File .env not found, reading configuration from ENV")
	}
	jmg = &JwtManager{}
	jmg.accessTokenBlackList = make(map[string]struct{})
	atl := os.Getenv("ACCESS_TOKEN_LIFE")
	rtl := os.Getenv("REFRESH_TOKEN_LIFE")
	atlInt, err := strconv.Atoi(atl)
	if err != nil {
		log.Fatal(err)
	}
	rtlInt, err := strconv.Atoi(rtl)
	if err != nil {
		log.Fatal(err)
	}
	jmg.AccessTokenLife = time.Duration(atlInt) * time.Minute
	jmg.RefreshTokenLife = time.Duration(rtlInt) * time.Hour
	go blackListCleaner()
}

func GetJWTManager() *JwtManager {
	return jmg
}

func (jmg *JwtManager) GenerateToken(ID, role string, exp time.Time) (string, error) {
	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	claims := jwt.MapClaims{}
	claims["user_id"] = ID
	claims["role"] = role
	claims["exp"] = exp.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("jwtmanager error: %s", err.Error())
	}
	return tokenString, nil
}

func (jmg *JwtManager) ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("token error")
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
}

func (jmg *JwtManager) AddTokenToBlackList(token string) {
	jmg.accessTokenBlackList[token] = struct{}{}
}

func (jmg *JwtManager) IsTokenInBlackList(token string) bool {
	if _, ok := jmg.accessTokenBlackList[token]; ok {
		return ok
	}
	return false
}

func (jmg *JwtManager) GenExp(typeToken string) time.Time {
	if typeToken == "access" {
		return time.Now().Add(jmg.AccessTokenLife)
	} else if typeToken == "refresh" {
		return time.Now().Add(jmg.RefreshTokenLife)
	} else {
		return time.Unix(0, 0)
	}
}

func blackListCleaner() {
	for range time.NewTicker(time.Duration(30 * time.Minute)).C {
		for key := range jmg.accessTokenBlackList {
			token, err := jmg.ParseToken(key)
			if err != nil {
				delete(jmg.accessTokenBlackList, key)
				continue
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if exp, ok := claims["exp"].(int64); ok {
					if exp < time.Now().Unix() {
						fmt.Println(key)
						delete(jmg.accessTokenBlackList, key)
					}
				}
			}
		}
	}
}
