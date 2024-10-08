package middleware

import (
	"fledge-restapi/internal/util"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization header"})
			c.Abort()
			return
		}

		// Remove 'Bearer ' prefix if present
		if len(tokenString) > 7 && strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = tokenString[7:]
		}

		claims, err := util.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}

func RateLimiter() gin.HandlerFunc {
	// Simple in-memory store for rate limiting
	type client struct {
		lastSeen time.Time
		count    int
	}

	var (
		mutex   sync.Mutex
		clients = make(map[string]*client)
	)

	return func(c *gin.Context) {
		mutex.Lock()
		defer mutex.Unlock()

		ip := c.ClientIP()
		now := time.Now()

		if cl, exists := clients[ip]; exists {
			if now.Sub(cl.lastSeen) > time.Minute {
				cl.count = 0
			}

			if cl.count >= 100 { // 100 requests per minute
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error": "Rate limit exceeded",
				})
				c.Abort()
				return
			}

			cl.count++
			cl.lastSeen = now
		} else {
			clients[ip] = &client{
				lastSeen: now,
				count:    1,
			}
		}

		c.Next()
	}
}

// Cors creates a new CORS middleware
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
