package main

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/go-ldap/ldap/v3"
	auth "github.com/korylprince/go-ad-auth/v3"
)

var stripRegexp = regexp.MustCompile("[^0-9]")

func pagedSearch(conn *auth.Conn, filter string, attrs []string) ([]*ldap.Entry, error) {
	search := ldap.NewSearchRequest(
		conn.Config.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.DerefAlways,
		0,
		0,
		false,
		filter,
		attrs,
		nil,
	)
	result, err := conn.Conn.SearchWithPaging(search, 1000)
	if err != nil {
		return nil, fmt.Errorf(`Search error "%s": %v`, filter, err)
	}

	return result.Entries, nil
}

//User represents a user
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

//GetADUsers returns all the users specified by the given config, or an error if one occurred
func GetADUsers(config *Config) ([]*User, error) {
	adConfig := &auth.Config{
		Server:   config.LDAPServer,
		Port:     config.LDAPPort,
		BaseDN:   config.LDAPBaseDN,
		Security: config.SecurityType(),
	}
	conn, err := adConfig.Connect()
	if err != nil {
		return nil, fmt.Errorf("Unable to connect: %v", err)
	}

	ok, err := conn.Bind(config.LDAPBindUser, config.LDAPBindPass)
	if err != nil {
		return nil, fmt.Errorf("Unable to bind: %v", err)
	}
	if !ok {
		return nil, errors.New("Unable to bind: authentication unsuccessful")
	}

	userEntries, err := pagedSearch(conn, config.LDAPUserFilter, []string{"displayName", config.LDAPMatchAttr})
	if err != nil {
		return nil, fmt.Errorf("Unable to search: %v", err)
	}

	users := make([]*User, 0, len(userEntries))
	for _, e := range userEntries {
		u := &User{
			ID:   e.GetAttributeValue(config.LDAPMatchAttr),
			Name: e.GetAttributeValue("displayName"),
		}
		if config.StripMatchAttr {
			u.ID = stripRegexp.ReplaceAllString(u.ID, "")
		}
		users = append(users, u)
	}

	return users, nil
}
