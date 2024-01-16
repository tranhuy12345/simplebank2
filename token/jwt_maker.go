package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey: secretKey}, nil

}
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayLoad(username, duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":       payload.ID,
		"username": payload.Username,
		"exp":      payload.ExpiredAt,
		"issue":    payload.IssuedAt,
	})

	tokenString, err := jwtToken.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("token is invalid")
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.Parse(token, keyFunc)
	fmt.Println(jwtToken)
	if err != nil {
		return nil, err
	}
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Token error")
	}
	fmt.Println(claims)

	username := claims["username"].(string)
	issueAt := int64(claims["issue"].(float64))
	exp := int64(claims["exp"].(float64))
	id := claims["ID"].(string)
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	//Validate ngày hết hạn
	currentTime := time.Now().Unix()
	if currentTime > exp {
		return nil, errors.New("Token is expired")
	}
	payload := &Payload{
		ID:        uuid,
		Username:  username,
		IssuedAt:  issueAt,
		ExpiredAt: exp,
	}
	return payload, nil

}
