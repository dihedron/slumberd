package authenticate

import (
	"log/slog"

	"github.com/dihedron/devws/command/api"
	"github.com/dihedron/devws/command/base"
)

type Authenticate struct {
	base.Command
	// Address is the addrss of the LDAP server, including the schema.
	Address string `short:"a" long:"address" description:"LDAP server address." default:"ldaps://ldap.example.com:636" required:"yes" env:"DEVWS_LDAP_ADDRESS"`
	// Username is the DN of the user used to bind to the server.
	Username string `short:"u" long:"username" description:"Bind user's Distinguished Name." default:"cn=service_account,ou=services,dc=example,dc=com" required:"yes" env:"DEVWS_LDAP_USERNAME"`
	// Password is the password of the user used to bind to the server.
	Password string `short:"p" long:"password" description:"Bind user's password." default:"<ServiceAccoutSecret>" required:"yes" env:"DEVWS_LDAP_PASSWORD"`
	// BaseDN is the base DN used to perform LDAP searches.
	BaseDN string `short:"b" long:"base-dn" description:"Base DN used for LDAP searches." default:"dc=example,dc=com" required:"yes" env:"DEVWS_LDAP_BASEDN"`
	// Args are the required username and password to validate
	Args struct {
		Username string `required:"yes"`
		Password string `required:"yes"`
	} `positional-args:"yes"`
}

func (cmd *Authenticate) Execute(args []string) error {
	slog.Debug("running authenticate command", "address", cmd.Address, "username", cmd.Username, "password", cmd.Password, "base DN", cmd.BaseDN, "args", args)

	// if len(args) < 2 {
	// 	slog.Error("invalid format: username and password must be provided")
	// 	return fmt.Errorf("invalid format: devws [options] <username> <password>")
	// }

	// create the authenticator
	authenticator, err := api.NewLDAPAuthenticator(cmd.Username, cmd.Password, cmd.Address, cmd.BaseDN)
	if err != nil {
		slog.Error("failed to create LDAP authenticator", "error", err)
		return err
	}
	defer authenticator.Close()

	// attempt the user authentication
	slog.Debug("attempting user authentication", "username", cmd.Args.Username, "password", cmd.Args.Password)
	if ok, err := authenticator.Authenticate(cmd.Args.Username, cmd.Args.Password); err == nil {
		if ok {
			slog.Debug("user successfully authenticated", "username", cmd.Args.Username)
		} else {
			slog.Debug("user authentication failed", "username", cmd.Args.Username)
		}
	} else {
		slog.Error("failed authenticationg user", "username", cmd.Args.Username, "error", err)
		return err
	}
	slog.Debug("authenticate command completed")
	return nil
}
