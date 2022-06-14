package apis

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/loivis/swaggerhub-go"
	"github.com/loivis/swaggerhub-go/handy"
	"github.com/peterbourgon/ff/v3/ffcli"
	"golang.org/x/sync/errgroup"
)

type updateCollaboratorsConfig struct {
	client *swaggerhub.Client

	owner string
	paths string
	// TODO: updating collaborators doesn't need a file as content of the API.
	// apis, which is names of the APIs, could be directly used instead.
	apis string

	// non-existing users/teams will added.
	// existing ones will be updated with new roles.
	// roles will be assigned to all users/teams.
	users string
	teams string
	roles string
}

func (c *updateCollaboratorsConfig) validate() error {
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

	c.roles = strings.ToUpper(c.roles)

	return nil
}

func updateCollaborators(client *swaggerhub.Client) *ffcli.Command {
	config := &updateCollaboratorsConfig{
		client: client,
	}

	fs := flag.NewFlagSet("swagger apis get", flag.ExitOnError)
	fs.StringVar(&config.owner, "owner", "", "owner of the api/apis, organization or user, case-sensitive")
	fs.StringVar(&config.paths, "paths", "", "paths of the api/apis, comma separated list of files or directories")
	fs.StringVar(&config.users, "users", "", "users for the api/apis, comma separated list of names")
	fs.StringVar(&config.teams, "teams", "", "teams for the api/apis, comma separated list of names")
	fs.StringVar(&config.roles, "roles", "VIEW", "comma separated list of roles for the api/apis, valid values: VIEW, COMMENT, EDIT")

	return &ffcli.Command{
		Name:    "update-collaborators",
		FlagSet: fs,
		Exec: func(ctx context.Context, args []string) error {
			return updateCollaboratorsExec(ctx, config)
		},
	}
}

func updateCollaboratorsExec(ctx context.Context, config *updateCollaboratorsConfig) error {
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

			// merge updates into existing collaboration
			// TODO: deal with collaboration.PendingMembers?
			roles := strings.Split(config.roles, ",")

			users := strings.Split(config.users, ",")
			for _, user := range users {
				user = strings.Trim(user, " ")
				if user == "" {
					continue
				}

				var found bool

				newMember := collaborationMembership{Name: user, Roles: roles}

				for i := range col.Members {
					if col.Members[i].Name == user {
						col.Members[i] = newMember
						found = true
						break
					}
				}

				if !found {
					col.Members = append(col.Members, newMember)
				}
			}

			teams := strings.Split(config.teams, ",")
			for _, t := range teams {
				t = strings.Trim(t, " ")
				if t == "" {
					continue
				}

				var found bool

				newTeam := collaborationMembership{Name: t, Roles: roles}

				for i := range col.Teams {
					if col.Teams[i].Name == t {
						col.Teams[i] = newTeam
						found = true
						break
					}
				}

				if !found {
					col.Teams = append(col.Teams, newTeam)
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
