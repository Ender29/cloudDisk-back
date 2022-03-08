package util

import (
	"github.com/golang-jwt/jwt"
	"time"
)

var jwtSecret = []byte("python")

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	//CurrentTime int64
	jwt.StandardClaims
}

func GenerateToken(username, password string, times time.Duration) (string, error) {
	claims := Claims{
		username,
		password,
		//time.Now().UnixNano(),
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(times).Unix(),
			Issuer:    "let's go",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

func ParseToken(token string) (*Claims, int8) {
	var status int8 = 0
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				status = 1
			} else {
				status = -1
			}
		}
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, status
		}
	}

	return nil, status
}

// RefreshToken 刷新token
//func RefreshToken(tokenStr string) (string, error) {
//	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
//		return jwtSecret, nil
//	})
//	if err != nil {
//		return "", err
//	}
//	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
//		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
//		return GenerateToken()
//	}
//	return "", TokenInvalid
//}
