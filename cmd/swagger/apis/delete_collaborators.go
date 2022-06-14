package apis

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/loivis/swaggerhub-go"
	"github.com/loivis/swaggerhub-go/handy"
	"github.com/peterbourgon/ff/v3/ffcli"
	"golang.org/x/sync/errgroup"
)

type deleteCollaboratorsConfig struct {
	client *swaggerhub.Client

	owner string
	paths string
	// TODO: deleting collaborators doesn't need a file as content of the API.
	// apis, which is names of the APIs, could be directly used instead.
	apis string

	users string
	teams string
}

func (c *deleteCollaboratorsConfig) validate() error {
	if c.client == nil {
		return errors.New("missing swaggerhub client")
	}

	if c.owner == "" {
		return errors.New("must specify -owner")
	}

	if c.paths == "" {
		return errors.New("must specify -path")
	}

	if c.users == "" && c.teams == "" {
		return errors.New("must specify -users or -teams, or both")
	}

	return nil
}

func deleteCollaborators(client *swaggerhub.Client) *ffcli.Command {
	config := &deleteCollaboratorsConfig{
		client: client,
	}

	fs := flag.NewFlagSet("swagger apis get", flag.ExitOnError)
	fs.StringVar(&config.owner, "owner", "", "owner of the api/apis, organization or user, case-sensitive")
	fs.StringVar(&config.paths, "paths", "", "paths of the api/apis, comma separated list of files or directories")
	fs.StringVar(&config.users, "users", "", "users for the api/apis, comma separated list of names")
	fs.StringVar(&config.teams, "teams", "", "teams for the api/apis, comma separated list of names")

	return &ffcli.Command{
		Name:    "delete-collaborators",
		FlagSet: fs,
		Exec: func(ctx context.Context, args []string) error {
			return deleteCollaboratorsExec(ctx, config)
		},
	}
}

func deleteCollaboratorsExec(ctx context.Context, config *deleteCollaboratorsConfig) error {
	if err := config.validate(); err != nil {
		return fmt.Errorf("error validating config: %v", err)
	}

	files, err := handy.APIFiles(strings.Split(config.paths, ","))
	if err != nil {
		return fmt.Errorf("error reading apis from paths(%s): %v", config.paths, err)
	}

	var g errgroup.Group

	for _, file := range files {
		file := file

		g.Go(func() error {
			// get current collaboration
			param := swaggerhub.APICollaborationGetParam{
				Owner: config.owner,
				API:   file.Name,
			}
			b, err := config.client.APICollaborationGet(param)
			if err != nil {
				return fmt.Errorf("error fetching collaboration for %s/%s: %v", config.owner, file.Name, err)
			}

			var col collaboration

			if err := json.Unmarshal(b, &col); err != nil {
				return fmt.Errorf("error unmarshalling response: %v", err)
			}

			// delete users/teams from existing collaboration
			// TODO: deal with collaboration.PendingMembers?

			users := strings.Split(config.users, ",")
			for _, user := range users {
				user = strings.Trim(user, " ")
				if user == "" {
					continue
				}

				for i := range col.Members {
					if col.Members[i].Name == user {
						col.Members = append(col.Members[:i], col.Members[i+1:]...)
						log.Printf("deleting user %q from %s", user, file.Name)
						break
					}
				}
			}

			teams := strings.Split(config.teams, ",")
			for _, t := range teams {
				t = strings.Trim(t, " ")
				if t == "" {
					continue
				}

				for i := range col.Teams {
					if col.Teams[i].Name == t {
						col.Members = append(col.Members[:i], col.Members[i+1:]...)
						log.Printf("deleting team %q from %s", t, file.Name)
						break
					}
				}
			}

			// push collaboration update
			b, err = json.Marshal(col)
			if err != nil {
				return fmt.Errorf("error marshalling collaboration: %v", err)
			}

			putParam := swaggerhub.APICollaborationPutParam{
				Owner: config.owner,
				API:   file.Name,
				Body:  b,
			}
			err = config.client.APICollaborationPut(putParam)
			if err != nil {
				return fmt.Errorf("error updating collaboration for %s/%s: %v", config.owner, file.Name, err)
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
