package api

import (
	"fmt"
	"log/slog"

	"github.com/go-ldap/ldap/v3"
)

type Authenticator interface {
	// Authenticate will return true if the user could be successfully
	// authenticated; false (with no error) if the user's credentials
	// are invalid; false (with an error) if the authenticator encountered
	// and internal processing error.
	Authenticate(username, password string) (bool, error)
	// Close can be used to perform cleanup operations.
	Close() error
}

// StaticAuthenticator authenticates users against an in memory, static map.
type StaticAuthenticator struct {
	accounts map[string]string
}

func NewStaticAuthenticator(options ...func(*StaticAuthenticator)) *StaticAuthenticator {
	auth := &StaticAuthenticator{
		accounts: map[string]string{},
	}
	for _, option := range options {
		option(auth)
	}
	return auth
}

func WithUser(username, password string) func(*StaticAuthenticator) {
	return func(a *StaticAuthenticator) {
		a.accounts[username] = password
	}
}

func (a *StaticAuthenticator) Authenticate(username, password string) (bool, error) {
	if pass, exists := a.accounts[username]; exists {
		slog.Debug("user successfuilly authenticated", "username", username, "password", password)
		return pass == password, nil
	}
	slog.Debug("error authenticating user", "username", username)
	return false, nil
}

func (a *StaticAuthenticator) Close() error {
	return nil
}

type LDAPAuthenticator struct {
	address    string
	account    string
	password   string
	basedn     string
	connection *ldap.Conn
	//filter     string
}

// NewLDAPAuthenticator initialises an LDAP authenticator using
// the given LDAP server address, service account and password;
// moreover is stores the BaseDN used for subsequent queries.
func NewLDAPAuthenticator(account, password, address, basedn string) (*LDAPAuthenticator, error) {

	slog.Debug("connecting to LDAP server", "address", address, "account", account, "password", password, "base DN", basedn)

	// connect to the LDAP server
	connection, err := ldap.DialURL(address)
	if err != nil {
		slog.Error("failed to connect to LDAP", "address", address, "error", err)
		return nil, fmt.Errorf("failed to connect to LDAP: %w", err)
	}
	slog.Debug("connected to LDAP server", "address", address)

	// bind with the service account to search the directory
	if err = connection.Bind(account, password); err != nil {
		slog.Error("failed to bind service account", "address", address, "account", account)
		return nil, fmt.Errorf("failed to bind service account: %w", err)
	}

	slog.Info("successfully connected to LDAP server")

	return &LDAPAuthenticator{
		address:    address,
		account:    account,
		password:   password,
		basedn:     basedn,
		connection: connection,
	}, nil
}

func (a *LDAPAuthenticator) Close() error {
	if a.connection != nil {
		return a.connection.Close()
	}
	return nil
}

func (a *LDAPAuthenticator) Authenticate(username, password string) (bool, error) {

	// search for the user's Distinguished Name (DN)
	search := ldap.NewSearchRequest(
		a.basedn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=person)(|(uid=%s)(sAMAccountName=%s)))", ldap.EscapeFilter(username), ldap.EscapeFilter(username)),
		[]string{"dn"}, // We only need to retrieve the DN, no other attributes
		nil,
	)

	result, err := a.connection.Search(search)
	if err != nil {
		slog.Error("failed to search for user", "username", username)
		return false, fmt.Errorf("failed to search for user: %w", err)
	}

	// Handle search results
	if len(result.Entries) == 0 {
		return false, fmt.Errorf("user not found")
	}
	if len(result.Entries) > 1 {
		return false, fmt.Errorf("multiple users found with the same username")
	}

	// Extract the user's exact DN from the search result
	dn := result.Entries[0].DN

	connection, err := ldap.DialURL(a.address)
	if err != nil {
		slog.Error("error connecting to LDAP server", "address", a.address, "error", err)
		return false, fmt.Errorf("failed to connect to LDAP: %w", err)
	}
	defer connection.Close()

	// Step 4: Re-Bind as the specific user to verify their password
	err = connection.Bind(dn, password)
	if err != nil {
		// If the error is LDAP Result Code 49 (Invalid Credentials), the password was wrong.
		// We return false, but no error, as this is an expected authentication failure.
		if ldap.IsErrorWithCode(err, ldap.LDAPResultInvalidCredentials) {
			return false, nil
		}
		// Any other error means the bind failed for a system reason (e.g., connection lost)
		return false, fmt.Errorf("failed to bind as user: %w", err)
	}

	// If the second bind succeeds, the credentials are valid!
	return true, nil
}

/*
	if err = connection.Bind(fmt.Sprintf("CN=%s,)", ldap.EscapeFilter(username)), password); err != nil {
		if ldap.IsErrorWithCode(err, ldap.LDAPResultInvalidCredentials) {
			slog.Warn("invalid credentials")
			return false //, nil Bad password
		}
		slog.Error("network error connectin to LDAP server", "address", a.address, "error", err)
		return false // , err Network/System error
	}
	slog.Debug("user successfully authenticated", "username", username)
	return true
}
*/
