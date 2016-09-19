package systems

import (
	"time"

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

// GenerateJWTToken will create new JWt Token
func (j *Jwt) GenerateJWTToken(userGUID string, phoneNo string, deviceUUID string) (*JwtToken, *ErrorData) {
	currentTimestamp := time.Now().Unix()

	// jti := md5.New()
	// jti.Write([]byte(deviceUUID))
	// jtiHash := string(jti.Sum(nil))

	// Create the token
	expired := time.Now().AddDate(0, 0, 7).Unix()
	// Create the Claims
	claims := CustomClaims{
		phoneNo,
		jwt.StandardClaims{
			Audience:  userGUID,
			ExpiresAt: expired,
			Id:        deviceUUID,
			IssuedAt:  currentTimestamp,
			Issuer:    "http://api.shoppermate.com",
			NotBefore: currentTimestamp,
			Subject:   userGUID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	config := Configs{}
	jwtSecret := config.Get("app.yaml", "jwt_token_secret", "secret")

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		ErrorMesg := &Error{}
		return nil, ErrorMesg.InternalServerError(err.Error(), FailedToGenerateToken)
	}

	tokenData := &JwtToken{}
	tokenData.Token = tokenString

	tokenExpiredDate := time.Unix(expired, 0)
	tokenData.Expired = tokenExpiredDate.String()
	return tokenData, nil
}
