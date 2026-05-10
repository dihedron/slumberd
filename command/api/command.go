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

	Address string `short:"a" long:"address" description:"Address to bind the API to." default:":8080"`
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
	Link   *Link  `json:"link,omitempty"`
	Links  *struct {
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
		Link:   &Link{Relation: "self", Href: base + "/" + id},
		Links: &struct {
			Stop     *Link `json:"stop,omitempty"`
			Start    *Link `json:"start,omitempty"`
			Restart  *Link `json:"restart,omitempty"`
			Pause    *Link `json:"pause,omitempty"`
			Unpause  *Link `json:"unpause,omitempty"`
			Shelve   *Link `json:"shelve,omitempty"`
			Unshelve *Link `json:"unshelve,omitempty"`
		}{
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

	// initialize the session store (using a cookie for simplicity here)
	// in production, use a strong, environment-variable-injected secret key.
	//store := cookie.NewStore(authenticationKey, encryptionKey)

	router.Use(Logger())
	router.Use(gin.Recovery())

	// register the session middleware FIRST so subsequent middlewares can use it
	router.Use(sessions.Sessions("api_session", cookie.NewStore([]byte("super-secret-key"))))

	// define an authenticator
	authenticator := &StaticAuthenticator{
		Accounts: map[string]string{
			"admin":     "QWERTY",
			"developer": "QWERTY",
		},
	}
	// group routes that require authentication
	router.Use(SessionAuthMiddleware("Developer Workstation", authenticator)) // Apply our custom middleware

	group := router.Group("/api/v1/vm")
	{
		// group.GET("/login", func(c *gin.Context) {
		// 	session := sessions.Default(c)
		// 	session.Set("user", "admin")
		// 	session.Save()
		// 	c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully"})
		// })
		group.GET("/logout", func(c *gin.Context) {
			session := sessions.Default(c)
			session.Clear()
			session.Save()
			c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
		})

		group.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, []*VM{
				NewVM("http://"+c.Request.Host+"/api/v1/vm", "vm001", "running"),
				NewVM("http://"+c.Request.Host+"/api/v1/vm", "vm002", "stopped"),
				NewVM("http://"+c.Request.Host+"/api/v1/vm", "vm003", "running"),
				NewVM("http://"+c.Request.Host+"/api/v1/vm", "vm004", "paused"),
				NewVM("http://"+c.Request.Host+"/api/v1/vm", "vm005", "shelved"),
			})
		})
		group.GET("/:vm", func(c *gin.Context) {
			c.JSON(http.StatusOK, NewVM("http://"+c.Request.Host+"/api/v1/vm", c.Param("vm"), "stopped"))
		})
	}

	slog.Info("API server running", "address", cmd.Address)
	err := router.Run(cmd.Address)
	if err != nil {
		slog.Error("API server failed", "error", err)
		return fmt.Errorf("API server failed: %w", err)
	}
	return nil
}
