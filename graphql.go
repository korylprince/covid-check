package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/korylprince/go-graphql-ws"
)

const gqlUpsertUsers = `
	mutation UpsertUsers($users: [user_insert_input!]!) {
	  insert_user(objects: $users, on_conflict: {constraint: user_pkey, update_columns: name}) {
		affected_rows
	  }
	}
`

//UpsertUsers upserts the given users to the server specified by the config
//returning the amount of users affeced, or an error if one occurred
func UpsertUsers(config *Config, users []*User) (int, error) {
	type response struct {
		InsertUser struct {
			AffectedRows int `json:"affected_rows"`
		} `json:"insert_user"`
	}

	headers := make(http.Header)
	if config.GraphQLAPISecret != "" {
		headers.Add("X-Hasura-Admin-Secret", config.GraphQLAPISecret)
	}

	var q = &graphql.MessagePayloadStart{
		Query: gqlUpsertUsers,
		Variables: map[string]interface{}{
			"users": users,
		},
	}

	conn, _, err := graphql.DefaultDialer.Dial(config.GraphQLEndpoint, headers, nil)
	if err != nil {
		return 0, fmt.Errorf("Unable to connect: %v", err)
	}

	payload, err := conn.Execute(context.Background(), q)
	if err != nil {
		return 0, fmt.Errorf("Unable to execute query: %v", err)
	}

	resp := new(response)

	if err = json.Unmarshal(payload.Data, resp); err != nil {
		return 0, fmt.Errorf("Unable to decode response: %v", err)
	}

	return resp.InsertUser.AffectedRows, nil
}
