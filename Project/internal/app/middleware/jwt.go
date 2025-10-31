package middleware

import (
	"net/http"
	"strings"

	"Project/internal/api/cache"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("secret_key")

func JWTMiddleware(rdb *cache.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// Получаем токен из Authorization или cookie
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else if cookie, err := c.Cookie("session_token"); err == nil {
			tokenString = cookie
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Не авторизован"})
			c.Abort()
			return
		}

		// Проверка черного списка
		isBlacklisted, err := rdb.IsTokenBlacklisted(tokenString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ошибка проверки токена"})
			c.Abort()
			return
		}
		if isBlacklisted {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Токен заблокирован"})
			c.Abort()
			return
		}

		// Парсим JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenMalformed
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Не авторизован"})
			c.Abort()
			return
		}

		// Сохраняем токен и данные пользователя в контекст
		c.Set("userToken", tokenString)
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if id, ok := claims["userID"].(float64); ok {
				c.Set("userID", uint(id))
			}
			if mod, ok := claims["isModerator"].(bool); ok {
				c.Set("isModerator", mod)
			} else {
				c.Set("isModerator", false)
			}
		}

		c.Next()
	}
}
