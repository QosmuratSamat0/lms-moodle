package deliveryhttp

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID   int64
	Name string
	Role string
}

const userKey = "user"

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := strings.TrimSpace(c.GetHeader("X-Role"))
		userIDRaw := strings.TrimSpace(c.GetHeader("X-User-Id"))
		userName := strings.TrimSpace(c.GetHeader("X-User-Name"))

		if role == "" || userIDRaw == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing auth headers"})
			return
		}

		userID, err := strconv.ParseInt(userIDRaw, 10, 64)
		if err != nil || userID <= 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
			return
		}

		c.Set(userKey, User{
			ID:   userID,
			Name: userName,
			Role: role,
		})
		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[strings.ToLower(strings.TrimSpace(role))] = struct{}{}
	}

	return func(c *gin.Context) {
		user, ok := UserFromContext(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		if _, ok := allowed[strings.ToLower(user.Role)]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}

func UserFromContext(c *gin.Context) (User, bool) {
	value, ok := c.Get(userKey)
	if !ok {
		return User{}, false
	}
	user, ok := value.(User)
	return user, ok
}
