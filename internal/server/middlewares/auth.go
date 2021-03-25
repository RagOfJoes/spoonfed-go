package middlewares

import (
	"net/http"
	"strings"

	"github.com/RagOfJoes/spoonfed-go/internal/orm"
	"github.com/RagOfJoes/spoonfed-go/internal/orm/services"
	"github.com/RagOfJoes/spoonfed-go/pkg/auth"
	"github.com/RagOfJoes/spoonfed-go/pkg/logger"
	"github.com/RagOfJoes/spoonfed-go/pkg/util"
	"github.com/gin-gonic/gin"
)

// Auth checks for an access token in the Header of the request
// with the key "Authorization" and with a value of Bearer ...TOKEN.
// If an access token was provided then run GetUser fn and provide its response to
// context for futher usage.
// If no access token was provided then continue with request execution.
func Auth(path string, db *orm.ORM) gin.HandlerFunc {
	logger.Infof("[Auth] attached to %s", path)
	return gin.HandlerFunc(func(c *gin.Context) {
		req := c.Request
		tokenValue := getTokenValue(req.Header)
		// 1. Check if valid token was provided
		if tokenValue == "" {
			c.Next()
			return
		}
		// 2. Retrieve OpenID client
		client, err := auth.GetClient()
		if err != nil {
			c.Next()
			return
		}
		// 3. GetUser with parsed access token
		user, err := client.GetUser(tokenValue)
		if err != nil || user == nil {
			c.Next()
			return
		}

		tx := db.DB.Begin()
		iUser, err := services.GetUserFromRagOfJoes(user, tx)
		if err != nil {
			c.Next()
			return
		}
		if err := tx.Commit().Error; err != nil {
			c.Next()
			return
		}
		// 4. Add to context then resume with the rest of the request
		util.AddToContext(c, util.ProjectContextKeys.User, iUser)
		c.Next()
	})
}

// Reads request header and validates "Authorization" header.
// If invalid then will return empty string("").
func getTokenValue(header http.Header) string {
	// 1. Check if Authorization header exists
	ah := header.Get(auth.AuthHeaderKey)
	if ah == "" {
		return ""
	}
	// 2. Split Authorization header
	split := strings.Split(ah, " ")
	if len(split) != 2 {
		return ""
	}
	// 3. Make sure token provided is a Bearer
	// token type
	tokenType, tokenValue := split[0], split[1]
	if !strings.EqualFold(tokenType, auth.ValidTokenType) {
		return ""
	}
	return tokenValue
}
