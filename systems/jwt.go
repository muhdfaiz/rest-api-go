package systems

import (
	"os"
	"time"

	"strconv"

	"github.com/dgrijalva/jwt-go"
)

type Jwt struct {
}

type JwtToken struct {
	Token   string `json:"token"`
	Expired string `json:"expired"`
}

type CustomClaims struct {
	PhoneNo string `json:"phone_no"`
	jwt.StandardClaims
}

// GenerateToken will create new JWt Token
func (j *Jwt) GenerateToken(userGUID string, phoneNo string, deviceUUID string, debugToken string) (*JwtToken, *ErrorData) {
	currentTimestamp := time.Now().UTC().Unix()

	// jti := md5.New()
	// jti.Write([]byte(deviceUUID))
	// jtiHash := string(jti.Sum(nil))

	expired := time.Now().UTC().AddDate(0, 0, 7).Unix()

	if debugToken != "" {
		debugTokenInInt, _ := strconv.Atoi(debugToken)
		minutes := time.Minute * time.Duration(debugTokenInInt)
		expired = time.Now().UTC().Add(minutes).Unix()
	}

	// Create the Claims
	claims := CustomClaims{
		phoneNo,
		jwt.StandardClaims{
			Audience:  userGUID,
			ExpiresAt: expired,
			Id:        deviceUUID,
			IssuedAt:  currentTimestamp,
			Issuer:    "http://api.shoppermate-api.com",
			NotBefore: currentTimestamp,
			Subject:   userGUID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := os.Getenv("JWT_TOKEN_SECRET")

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		Error := &Error{}
		return nil, Error.InternalServerError(err.Error(), FailedToGenerateToken)
	}

	tokenData := &JwtToken{}
	tokenData.Token = tokenString

	tokenExpiredDate := time.Unix(expired, 0).UTC().Format(time.RFC3339)
	tokenData.Expired = tokenExpiredDate

	return tokenData, nil
}
