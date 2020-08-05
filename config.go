package main

import (
	"log"
	"strings"
	"time"

	auth "github.com/korylprince/go-ad-auth/v3"
)

//Config is the configuration for the application
type Config struct {
	LDAPServer       string        `required:"true"`
	LDAPPort         int           `default:"389"`
	LDAPBaseDN       string        `required:"true"`
	LDAPSecurity     string        `default:"none"`
	LDAPBindUser     string        `required:"true"`
	LDAPBindPass     string        `required:"true"`
	LDAPUserFilter   string        `default:"(objectClass=user)"`
	LDAPMatchAttr    string        `default:"employeeID"`
	StripMatchAttr   bool          `default:"true"` //strip all but numbers from ID
	SyncInterval     time.Duration `default:"30m"`
	GraphQLEndpoint  string        `required:"true"`
	GraphQLAPISecret string
}

//SecurityType returns the auth.SecurityType for the config
func (c *Config) SecurityType() auth.SecurityType {
	switch strings.ToLower(c.LDAPSecurity) {
	case "", "none":
		return auth.SecurityNone
	case "tls":
		return auth.SecurityTLS
	case "starttls":
		return auth.SecurityStartTLS
	default:
		log.Fatalln("ERROR: Invalid LDAPSECURITY:", c.LDAPSecurity)
	}
	panic("unreachable")
}
