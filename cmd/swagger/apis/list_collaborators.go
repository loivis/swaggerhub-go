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

type listCollaboratorsConfig struct {
	client *swaggerhub.Client

	owner string
	paths string
	// TODO: deleting collaborators doesn't need a file as content of the API.
	// apis, which is names of the APIs, could be directly used instead.
	apis string
}

func (c *listCollaboratorsConfig) validate() error {
	if c.client == nil {
		return errors.New("missing swaggerhub client")
	}

	if c.owner == "" {
		return errors.New("must specify -owner")
	}

	if c.paths == "" {
		return errors.New("must specify -path")
	}

	return nil
}

func listCollaborators(client *swaggerhub.Client) *ffcli.Command {
	config := &listCollaboratorsConfig{
		client: client,
	}

	fs := flag.NewFlagSet("swagger apis get", flag.ExitOnError)
	fs.StringVar(&config.owner, "owner", "", "owner of the api/apis, organization or user, case-sensitive")
	fs.StringVar(&config.paths, "paths", "", "paths of the api/apis, comma separated list of files or directories")

	return &ffcli.Command{
		Name:    "list-collaborators",
		FlagSet: fs,
		Exec: func(ctx context.Context, args []string) error {
			return listCollaboratorsExec(ctx, config)
		},
	}
}

func listCollaboratorsExec(ctx context.Context, config *listCollaboratorsConfig) error {
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

			log.Printf("%s: %+v", file.Name, col)

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
