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

type setDefaultConfig struct {
	client *swaggerhub.Client

	owner string
	paths string

	version string
}

func (c *setDefaultConfig) validate() error {
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

func setDefault(client *swaggerhub.Client) *ffcli.Command {
	config := &setDefaultConfig{
		client: client,
	}

	fs := flag.NewFlagSet("swagger apis get", flag.ExitOnError)
	fs.StringVar(&config.owner, "owner", "", "owner of the apis, organization or user, case-sensitive")
	fs.StringVar(&config.paths, "paths", "", "paths of api/apis, comma separated list of files or directories")
	fs.StringVar(&config.version, "version", "", "version of the api, info.version in the specification will be used if omitted")

	return &ffcli.Command{
		Name:    "set-default",
		FlagSet: fs,
		Exec: func(ctx context.Context, args []string) error {
			return setDefaultExec(ctx, config)
		},
	}
}

func setDefaultExec(ctx context.Context, config *setDefaultConfig) error {
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
			param := swaggerhub.APISettingsDefaultPutParam{
				Owner:   config.owner,
				API:     file.Name,
				Version: version,
			}
			if err := config.client.APISettingsDefaultPut(param); err != nil {
				return fmt.Errorf("error setting default for %s/%s to %s: %v", config.owner, file.Name, config.version, err)
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
