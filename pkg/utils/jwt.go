package utils

import (
    "github.com/golang-jwt/jwt/v5"
    "time"
)

func GenerateToken(userID int64, secret string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 72).Unix(),
    })
    return token.SignedString([]byte(secret))
}

func ParseToken(tokenString, secret string) (int64, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    })
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return int64(claims["user_id"].(float64)), nil
    }
    return 0, err
}
