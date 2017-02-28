package middlewares

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"os"

	"bitbucket.org/cliqers/shoppermate-api/systems"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type CustomClaims struct {
	PhoneNo string `json:"phone_no"`
	jwt.StandardClaims
}

func Auth(DB *gorm.DB) gin.HandlerFunc {
	jwtSecret := os.Getenv("JWT_TOKEN_SECRET")

	return func(c *gin.Context) {
		authorizationHeader := c.Request.Header["Authorization"]

		if authorizationHeader == nil {
			Error := &systems.Error{}
			c.JSON(http.StatusUnauthorized, Error.GenericError(strconv.Itoa(http.StatusUnauthorized), systems.TokenNotValid,
				systems.TitleErrorTokenNotValid, "", systems.ErrorTokenNotValid))
			c.Abort()
			return
		}

		splitAuthorizationHeader := strings.SplitN(authorizationHeader[0], " ", 2)

		if len(splitAuthorizationHeader) != 2 {
			Error := &systems.Error{}
			c.JSON(http.StatusUnauthorized, Error.GenericError(strconv.Itoa(http.StatusUnauthorized), systems.TokenNotValid,
				systems.TitleErrorTokenNotValid, "", systems.ErrorTokenNotValid))
			c.Abort()
			return
		}

		tokenString := splitAuthorizationHeader[1]

		// Parse takes the token string and a function for looking up the key. The latter is especially
		// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
		// head of the token to identify which key to use, but the parsed token (head and claims) is provided
		// to the callback, providing flexibility.
		token, _ := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(jwtSecret), nil
		})

		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			type Device struct {
				ID           int        `json:"id"`
				GUID         string     `json:"guid"`
				UserGUID     *string    `json:"user_guid"`
				UUID         string     `json:"uuid"`
				Os           string     `json:"os"`
				Model        string     `json:"model"`
				PushToken    string     `json:"push_token"`
				AppVersion   string     `json:"app_version"`
				TokenExpired int        `json:"token_expired"`
				CreatedAt    time.Time  `json:"created_at"`
				UpdatedAt    time.Time  `json:"updated_at"`
				DeletedAt    *time.Time `json:"deleted_at"`
			}

			device := &Device{}

			result := DB.Where("uuid = ? AND user_guid = ?", claims.Id, claims.Subject).Find(&device)

			if result.RowsAffected == 0 {
				c.JSON(http.StatusUnauthorized, Error.GenericError(strconv.Itoa(http.StatusUnauthorized), systems.TokenNotValid,
					systems.TitleErrorTokenNotValid, "", systems.ErrorTokenNotValid))
				c.Abort()
				return
			}

			tokenData := make(map[string]string)
			tokenData["device_uuid"] = claims.Id
			tokenData["user_guid"] = claims.Subject
			tokenData["user_phone_no"] = claims.PhoneNo

			c.Set("Token", tokenData)
			c.Next()
			return
		}

		c.JSON(http.StatusUnauthorized, Error.GenericError(strconv.Itoa(http.StatusUnauthorized), systems.TokenNotValid,
			systems.TitleErrorTokenNotValid, "", systems.ErrorTokenNotValid))
		c.Abort()
		return

	}
}
