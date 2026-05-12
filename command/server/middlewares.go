package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Logger logs requests in a structured format.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		slog.Info("request performed", "method", c.Request.Method, "path", path, "query", query, "status", c.Writer.Status(), "latency", time.Since(start), "client IP", c.ClientIP(), "body size", c.Writer.Size())

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				slog.Error("request error", "error", err.Error())
			}
		}
	}
}

// SessionAuthMiddleware handles the combined Session + Basic Auth logic
func SessionAuthMiddleware(realm string, authenticator Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Debug("sesion manager middleware - START")
		session := sessions.Default(c)
		user := session.Get("username")

		// check if the user already has a valid session
		if user != nil {
			slog.Debug("valid session for user", "username", user)
			// session is valid, proceed to the API handler
			c.Next()
			return
		}

		// no valid session, check for Basic Authentication headers
		username, password, hasAuth := c.Request.BasicAuth()
		if hasAuth {
			if ok, _ := authenticator.Authenticate(username, password); ok {
				// Basic Auth is valid, create a session for future requests.
				session.Set("username", username)
				if err := session.Save(); err != nil {
					slog.Error("failed to save session", "error", err)
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
					return
				}
				// proceed to the API handler
				c.Next()
				return
			}
		}

		// no valid session and no/invalid Basic Auth; challenge the client.
		// the WWW-Authenticate header triggers the browser's native login prompt
		// or tells API clients (like curl/Postman) to provide Basic Auth.
		// c.Header("WWW-Authenticate", `Basic realm="`+realm+`"`)
		// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		// 	"error": "Unauthorized. Please provide valid credentials.",
		// })
		c.Redirect(http.StatusFound, "/login")
	}
}
