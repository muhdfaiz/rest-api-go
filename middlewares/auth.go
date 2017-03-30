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

// CustomClaims define the payload structure of JWT Token.
// Add additional payload called phone number.
// By default, JWT Standard Claims only contain payload below:
// - Audience
// - ExpiresAt
// - Id
// - IssuedAt
// - Issuer
// - NotBefore
// - Subject
type CustomClaims struct {
	PhoneNo string `json:"phone_no"`
	jwt.StandardClaims
}

// Auth Middleware. Use to protect API endpoint that require authorization.
// It will validate JWT token valid or not.
func Auth(DB *gorm.DB) gin.HandlerFunc {
	jwtSecret := os.Getenv("JWT_TOKEN_SECRET")

	return func(c *gin.Context) {
		// Retrieve JWT token from `Authorization` key in request header.
		// Example valid authorization header:
		// Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwaG9uZV9ubyI6IjYwMTc0ODYyMTI3IiwiYXVkIjoiOGMyZ
		authorizationHeader := c.Request.Header["Authorization"]

		Error := &systems.Error{}

		// Return an error and abort the request if JWT token not exist in request header.
		if authorizationHeader == nil {
			c.JSON(http.StatusUnauthorized, Error.GenericError(strconv.Itoa(http.StatusUnauthorized), systems.TokenNotValid,
				systems.TitleErrorTokenNotValid, "", systems.ErrorTokenNotValid))

			c.Abort()

			return
		}

		// Split JWT token by empty space. Right now JWT token contain `Bearer` keyword in front of it.
		// Example: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwaG9uZV9ubyI6IjYwMTc0ODYyMTI3IiwiYXVkIjoiOGMyZ
		splitAuthorizationHeader := strings.SplitN(authorizationHeader[0], " ", 2)

		// Check if JWT token was specified in correct way. Must have `Bearer`keyword in front following with space and then the JWT token.
		// Return an error if JWT token not specified using correct way.
		if len(splitAuthorizationHeader) != 2 {
			c.JSON(http.StatusUnauthorized, Error.GenericError(strconv.Itoa(http.StatusUnauthorized), systems.TokenNotValid,
				systems.TitleErrorTokenNotValid, "", systems.ErrorTokenNotValid))

			c.Abort()

			return
		}

		// Retrieve JWT token only.
		tokenString := splitAuthorizationHeader[1]

		token, _ := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Check if JWT token using same algorithm with algorithm use during generate JWT token.
			// In this API, HMAC algorithm has been used to generate JWT token.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(jwtSecret), nil
		})

		// Check if JWT token valid or not.
		// If JWT token not valid return an error and abort the request.
		// If JWT Token valid, check device exist or not in database using device uuid and user guid in the token payload.
		// If device exist, continue request.
		// If device not exist, return an error and abort the request.
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

			// Set token data into the context.
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
