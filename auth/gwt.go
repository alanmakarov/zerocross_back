package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"log"
	"time"
)

// SecretKey — секретный ключ для подписи токенов
var SecretKey = []byte("your-secret-key") // Вы можете изменить этот ключ на более сложный

// Claims — структура, которая будет содержать информацию о пользователе
type Claims struct {
	Username string `json:"username"`
	UserID   int
	jwt.RegisteredClaims
}

// GenerateToken — функция для генерации JWT токена
func GenerateToken(username string, userId int) (string, error) {
	claims := Claims{
		Username: username,
		UserID:   userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Токен будет действителен 24 часа
			Issuer:    "fiber_api_v1",                                     // Указание источника токена
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(SecretKey)
	if err != nil {
		log.Println("Error signing token:", err)
		return "", err
	}

	return signedToken, nil
}
