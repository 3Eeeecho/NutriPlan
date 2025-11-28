package jwt

import (
	"NutriPlan/internal/config"
	"NutriPlan/pkg/xerr"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims 定义JWT的声明结构
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func GenerateToken(userID uint, username string) (string, error) {
	// 设置过期时间为7天
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	// 创建Claims
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// 使用HS256算法和密钥创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.AppConfig.Server.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("%w: %w", xerr.ErrTokenGeneration, err)
	}

	return tokenString, nil
}

// ValidateToken 验证JWT token并返回Claims
func ValidateToken(tokenString string) (*Claims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", xerr.ErrUnexpectedSigningMethod, token.Header["alg"])
		}
		return []byte(config.AppConfig.Server.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %w", xerr.ErrTokenParsing, err)
	}

	// 验证token是否有效
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, xerr.ErrInvalidToken
}

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{
				"error": "缺少Authorization请求头",
			})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{
				"error": "Authorization格式错误，应为: Bearer <token>",
			})
			c.Abort()
			return
		}

		// 提取token
		tokenString := parts[1]

		// 验证token
		claims, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{
				"error": "token验证失败: " + err.Error(),
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中，供后续处理器使用
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		// 继续处理请求
		c.Next()
	}
}
