package middlewares

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"bitbucket.org/shoppermate-api/application/v1"
	"bitbucket.org/shoppermate-api/systems"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type CustomClaims struct {
	PhoneNo string `json:"phone_no"`
	jwt.StandardClaims
}

func Auth() gin.HandlerFunc {
	jwtSecret := Config.Get("app.yaml", "jwt_token_secret", "secret")

	return func(c *gin.Context) {
		db := c.MustGet("DB").(*gorm.DB)
		authorizationHeader := c.Request.Header["Authorization"]

		if authorizationHeader == nil {
			ErrorMesg := &systems.Error{}
			c.JSON(http.StatusBadRequest, ErrorMesg.GenericError(strconv.Itoa(http.StatusBadRequest), systems.TokenNotValid,
				systems.TitleErrorTokenNotValid, "", systems.ErrorTokenNotValid))
			c.Abort()
			return
		}

		splitAuthorizationHeader := strings.SplitN(authorizationHeader[0], " ", 2)

		if len(splitAuthorizationHeader) != 2 {
			ErrorMesg := &systems.Error{}
			c.JSON(http.StatusBadRequest, ErrorMesg.GenericError(strconv.Itoa(http.StatusBadRequest), systems.TokenNotValid,
				systems.TitleErrorTokenNotValid, "", systems.ErrorTokenNotValid))
			c.Abort()
			return
		}

		tokenString := splitAuthorizationHeader[1]

		// Parse takes the token string and a function for looking up the key. The latter is especially
		// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
		// head of the token to identify which key to use, but the parsed token (head and claims) is provided
		// to the callback, providing flexibility.
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(jwtSecret), nil
		})

		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			result := db.Where("uuid = ? AND user_guid = ?", claims.Id, claims.Subject).Find(&v1.Device{})

			if result.RowsAffected == 0 {
				c.JSON(http.StatusBadRequest, ErrorMesg.GenericError(strconv.Itoa(http.StatusBadRequest), systems.TokenNotValid,
					systems.TitleErrorTokenNotValid, "", systems.ErrorTokenNotValid))
				c.Abort()
				return
			}

			tokenData := make(map[string]string)
			tokenData["device_uuid"] = claims.Id
			tokenData["user_guid"] = claims.Subject
			c.Set("Token", tokenData)
			c.Next()
			return
		}

		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, ErrorMesg.GenericError(strconv.Itoa(http.StatusBadRequest), systems.TokenNotValid,
			systems.TitleErrorTokenNotValid, "", systems.ErrorTokenNotValid))
		c.Abort()

		return

	}
}
