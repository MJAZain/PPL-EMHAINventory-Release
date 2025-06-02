package middleware

import (
	"encoding/json"
	"fmt"
	"go-gin-auth/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("PPL-K4-2025")

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		// Lanjutkan ke handler
		// Ambil claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.Respond(c, http.StatusUnauthorized, "Unauthorized", "Invalid token claims", nil)
			c.Abort()
			return
		}
		if ok {
			jsonBytes, err := json.MarshalIndent(claims, "", "  ")
			if err != nil {
				fmt.Println("Error encoding claims to JSON:", err)
			} else {
				fmt.Println("JWT Claims (JSON):")
				fmt.Println(string(jsonBytes))
			}
		}
		c.Set("user_id", claims["user_id"])
		c.Set("full_name", claims["full_name"])
		c.Next()
	}
}

func AuthAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil token dari header Authorization: Bearer <token>
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Respond(c, http.StatusUnauthorized, "Unauthorized", "Missing Authorization header", nil)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			utils.Respond(c, http.StatusUnauthorized, "Unauthorized", "Invalid token", nil)
			c.Abort()
			return
		}

		// Ambil claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.Respond(c, http.StatusUnauthorized, "Unauthorized", "Invalid token claims", nil)
			c.Abort()
			return
		}
		if ok {
			jsonBytes, err := json.MarshalIndent(claims, "", "  ")
			if err != nil {
				fmt.Println("Error encoding claims to JSON:", err)
			} else {
				fmt.Println("JWT Claims (JSON):")
				fmt.Println(string(jsonBytes))
			}
		}

		// Cek role
		if role, ok := claims["role"].(string); !ok || role != "admin" {
			utils.Respond(c, http.StatusForbidden, "Forbidden", "Admin access only", nil)
			c.Abort()
			return
		}

		// Lanjutkan ke handler
		c.Set("user_id", claims["user_id"])
		c.Set("full_name", claims["full_name"])
		c.Next()
	}
}
