package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/dihedron/devws/command/base"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type API struct {
	base.Command

	Address string `short:"a" long:"address" description:"Address to bind the API to." default:":3000"`
}

type Link struct {
	Relation string `json:"rel"`
	Href     string `json:"href"`
}

func (l Link) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"rel":  l.Relation,
		"href": l.Href,
	})
}

type VM struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Links  *struct {
		Self     *Link `json:"self,omitempty"`
		Stop     *Link `json:"stop,omitempty"`
		Start    *Link `json:"start,omitempty"`
		Restart  *Link `json:"restart,omitempty"`
		Pause    *Link `json:"pause,omitempty"`
		Unpause  *Link `json:"unpause,omitempty"`
		Shelve   *Link `json:"shelve,omitempty"`
		Unshelve *Link `json:"unshelve,omitempty"`
	} `json:"links,omitempty"`
}

func NewVM(base string, id string, status string) *VM {
	return &VM{
		ID:     id,
		Status: status,
		//Link:   &Link{Relation: "self", Href: base + "/" + id},
		Links: &struct {
			Self     *Link `json:"self,omitempty"`
			Stop     *Link `json:"stop,omitempty"`
			Start    *Link `json:"start,omitempty"`
			Restart  *Link `json:"restart,omitempty"`
			Pause    *Link `json:"pause,omitempty"`
			Unpause  *Link `json:"unpause,omitempty"`
			Shelve   *Link `json:"shelve,omitempty"`
			Unshelve *Link `json:"unshelve,omitempty"`
		}{
			Self:     &Link{Relation: "self", Href: base + "/" + id},
			Stop:     &Link{Relation: "stop", Href: base + "/" + id + "/stop"},
			Start:    &Link{Relation: "start", Href: base + "/" + id + "/start"},
			Restart:  &Link{Relation: "restart", Href: base + "/" + id + "/restart"},
			Pause:    &Link{Relation: "pause", Href: base + "/" + id + "/pause"},
			Unpause:  &Link{Relation: "unpause", Href: base + "/" + id + "/unpause"},
			Shelve:   &Link{Relation: "shelve", Href: base + "/" + id + "/shelve"},
			Unshelve: &Link{Relation: "unshelve", Href: base + "/" + id + "/unshelve"},
		},
	}
}

func (cmd *API) Execute(args []string) error {
	slog.Info("starting API server", "address", cmd.Address)

	router := gin.New()
	router.SetTrustedProxies(nil)

	// generate a session key from random bytes
	// this is used to secure the session cookie
	// authenticationKey := make([]byte, 32)
	// rand.Read(authenticationKey)
	// encryptionKey := make([]byte, 32)
	// rand.Read(encryptionKey)
	// store := cookie.NewStore(authenticationKey, encryptionKey)

	router.Use(
		Logger(),
		gin.Recovery(),
		sessions.Sessions("api_session", cookie.NewStore([]byte("super-secret-key"))),
	)

	// define an authenticator
	authenticator := NewStaticAuthenticator(
		WithUser("admin", "QWERTY"),
		WithUser("developer", "QWERTY"),
	)

	// the /login and /logout routes do not need authentication
	unauthenticated := router.Group("")
	{
		unauthenticated.StaticFile("/favicon.ico", "./command/api/assets/favicon.ico")
		unauthenticated.StaticFile("/devws.png", "./command/api/assets/devws.png")
		unauthenticated.StaticFile("/login", "./command/api/assets/login.html")
		unauthenticated.StaticFile("/style.css", "./command/api/assets/style.css")
		unauthenticated.StaticFile("/background.jpg", "./command/api/assets/background.jpg")
		unauthenticated.GET("/", func(c *gin.Context) {
			session := sessions.Default(c)
			if username := session.Get("username"); username != nil {
				slog.Debug("user already logged in, redirecting to main page...")
				c.Redirect(http.StatusFound, "/api/v1/vm/")
			} else {
				slog.Debug("user not logged in yet, redirecting to login page")
				c.Redirect(http.StatusFound, "/login")
			}
		})
		unauthenticated.POST("/login", func(c *gin.Context) {
			username := c.PostForm("username")
			password := c.PostForm("password")
			slog.Debug("logging out user first...", "username", username)
			session := sessions.Default(c)
			if u := session.Get("username"); u == username {
				slog.Debug("user already logged in, redirectong to main page")
				c.Redirect(http.StatusFound, "/api/v1/vm")
			} else {
				slog.Debug("logging in user...", "username", username, "password", password)
				if ok, err := authenticator.Authenticate(c, username, password); ok {
					slog.Info("user successfully logged in", "username", username)
					session.Set("username", username)
					session.Save()
					c.Redirect(http.StatusFound, "/api/v1/vm")
					return
				} else {
					slog.Error("failed to autheticate user", "username", username, "error", err)
				}
				c.Redirect(http.StatusFound, "/login")
			}

		})
		unauthenticated.GET("/logout", func(c *gin.Context) {
			session := sessions.Default(c)
			if username := session.Get("username"); username != nil {
				slog.Debug("logging out user...", "username", username)
				session.Clear()
				session.Save()
			}
			c.Redirect(http.StatusFound, "/api/v1/vm")
		})
	}

	authenticated := router.Group("/api/v1/vm", SessionAuthMiddleware("Developer Workstations Realm", authenticator))
	{

		authenticated.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, []*VM{
				NewVM("http://"+c.Request.Host+"/api/v1/vm", "vm001", "running"),
				NewVM("http://"+c.Request.Host+"/api/v1/vm", "vm002", "stopped"),
				NewVM("http://"+c.Request.Host+"/api/v1/vm", "vm003", "running"),
				NewVM("http://"+c.Request.Host+"/api/v1/vm", "vm004", "paused"),
				NewVM("http://"+c.Request.Host+"/api/v1/vm", "vm005", "shelved"),
			})
		})
		authenticated.GET("/:vm", func(c *gin.Context) {
			c.JSON(http.StatusOK, NewVM("http://"+c.Request.Host+"/api/v1/vm", c.Param("vm"), "stopped"))
		})
	}

	// /login
	// https://github.com/puikinsh/login-forms/tree/main/forms/glassmorphism
	// https://github.com/puikinsh/login-forms/tree/main/forms/material

	slog.Info("API server running", "address", cmd.Address)
	err := router.Run(cmd.Address)
	if err != nil {
		slog.Error("API server failed", "error", err)
		return fmt.Errorf("API server failed: %w", err)
	}
	return nil
}
