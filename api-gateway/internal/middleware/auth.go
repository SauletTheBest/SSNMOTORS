package middleware

import (
    "net/http"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

// Секретный ключ для подписи JWT
var jwtSecret = []byte("your-secret-key")

// Claims определяет структуру токена
type Claims struct {
    UserID string `json:"user_id"`
    jwt.RegisteredClaims
}

// AuthMiddleware проверяет JWT токен в заголовке Authorization
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Получаем заголовок Authorization
        auth := c.GetHeader("Authorization")
        if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
            return
        }

        // Извлекаем токен из заголовка
        tokenStr := strings.TrimPrefix(auth, "Bearer ")
        token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })

        // Проверяем валидность токена
        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }

        // Извлекаем данные пользователя из токена
        claims := token.Claims.(*Claims)
        c.Set("userID", claims.UserID)
        c.Next()
    }
}

// GenerateToken создает новый JWT токен для пользователя
func GenerateToken(userID string) (string, error) {
    claims := &Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Токен действителен 24 часа
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}