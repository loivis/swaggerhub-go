package apis

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/loivis/swaggerhub-go"
	"github.com/loivis/swaggerhub-go/handy"
	"github.com/peterbourgon/ff/v3/ffcli"
	"golang.org/x/sync/errgroup"
)

type updateConfig struct {
	client *swaggerhub.Client

	owner string
	paths string

	version   string
	isPrivate bool
}

func (c *updateConfig) validate() error {
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

func update(client *swaggerhub.Client) *ffcli.Command {
	config := &updateConfig{
		client: client,
	}

	fs := flag.NewFlagSet("swagger apis get", flag.ExitOnError)
	fs.StringVar(&config.owner, "owner", "", "owner of the apis, organization or user, case-sensitive")
	fs.StringVar(&config.paths, "paths", "", "paths of api/apis, comma separated list of files or directories")
	fs.StringVar(&config.version, "version", "", "version of the api, info.version in the specification will be used if omitted")
	fs.BoolVar(&config.isPrivate, "is-private", false, "whether to make the API private (true) or public (false)")

	return &ffcli.Command{
		Name:    "update",
		FlagSet: fs,
		Exec: func(ctx context.Context, args []string) error {
			return updateExec(ctx, config)
		},
	}
}

func updateExec(ctx context.Context, config *updateConfig) error {
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

		version := file.Version
		if config.version != "" && config.version != file.Version {
			return fmt.Errorf("error version mismatch: %q in flag vs %q from file", config.version, file.Version)
		}

		g.Go(func() error {
			param := swaggerhub.APIPostParam{
				Owner:       config.owner,
				API:         file.Name,
				Version:     version,
				IsPrivate:   config.isPrivate,
				ContentType: swaggerhub.ContentType{Request: "application/" + file.Format},
				Body:        file.File,
			}
			if err := config.client.APIPost(param); err != nil {
				return fmt.Errorf("error updating %s/%s/%s: %v", config.owner, file.Name, config.version, err)
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
