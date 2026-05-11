package authenticate

import (
	"log/slog"

	"github.com/dihedron/devws/command/api"
	"github.com/dihedron/devws/command/base"
)

type Authenticate struct {
	base.Command
	// Address is the addrss of the LDAP server, including the schema.
	Address string `short:"a" long:"address" description:"LDAP server address." default:"ldaps://ldap.example.com:636" required:"true" env:"DEVWS_LDAP_ADDRESS"`
	// Username is the DN of the user used to bind to the server.
	Username string `short:"u" long:"username" description:"Bind user's Distinguished Name." default:"cn=service_account,ou=services,dc=example,dc=com" required:"true" env:"DEVWS_LDAP_USERNAME"`
	// Password is the password of the user used to bind to the server.
	Password string `short:"p" long:"password" description:"Bind user's password." default:"<ServiceAccoutSecret>" required:"true" env:"DEVWS_LDAP_PASSWORD"`
	// BaseDN is the base DN used to perform LDAP searches.
	BaseDN string `short:"b" long:"base-dn" description:"Base DN used for LDAP searches." default:"dc=example,dc=com"required:"true" env:"DEVWS_LDAP_BASEDN"`
}

func (cmd *Authenticate) Execute(args []string) error {
	slog.Debug("running authenticate command", "address", cmd.Address, "username", cmd.Username, "password", cmd.Password, "base DN", cmd.BaseDN, "args", args)

	auth, err := api.NewLDAPAuthenticator(cmd.Username, cmd.Password, cmd.Address, cmd.BaseDN)
	if err != nil {
		slog.Error("failed to create LDAP authenticator", "error", err)
		return err
	}
	defer auth.Close()
	//auth.Authenticate()

	slog.Debug("authenticate command completed")
	return nil
}
