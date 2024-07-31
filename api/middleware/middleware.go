package middleware

import (
	"api-geteway/api/token"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func IsAuthenticated() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Cookie'dan tokenni olish
		tokenString, err := ctx.Cookie("access_token")
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Error": "Authorization cookie not found",
			})
			ctx.Abort()
			return
		}

		// JWT tokenni tekshirish va tasdiqlash
		claims, err := token.ExtractClaims(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Error": "Invalid token",
			})
			ctx.Abort()
			return
		}

		// Foydalanuvchi ma'lumotlarini context ga qo'shish
		ctx.Set("claims", claims)

		ctx.Next()
	}
}

func Authorize(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, exists := ctx.Get("claims")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Error": "Unauthorized",
			})
			ctx.Abort()
			return
		}

		userClaims := claims.(*token.Claims)
		fmt.Println(userClaims.Role, ctx.FullPath(), ctx.Request.Method)
		ok, err := enforcer.Enforce(userClaims.Role, ctx.FullPath(), ctx.Request.Method)
		if err != nil {
			fmt.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Error": "Internal server error",
				"Err": err.Error(),
			})
			ctx.Abort()
			return
		}

		if !ok {
			ctx.JSON(http.StatusForbidden, gin.H{
				"Error": "Forbidden",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func LogMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("Request received",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
		)

		c.Next()

		logger.Info("Response sent",
			slog.Int("status", c.Writer.Status()),
		)
	}
}
