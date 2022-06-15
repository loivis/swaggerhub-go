package domains

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

type getConfig struct {
	client *swaggerhub.Client

	owner string
	paths string

	version string
	format  string
}

func (c *getConfig) validate() error {
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

func sync(client *swaggerhub.Client) *ffcli.Command {
	config := &getConfig{
		client: client,
	}

	fs := flag.NewFlagSet("swagger apis get", flag.ExitOnError)
	fs.StringVar(&config.owner, "owner", "", "owner of the apis, organization or user, case-sensitive")
	fs.StringVar(&config.paths, "paths", "", "paths of api/apis, comma separated list of files or directories")
	fs.StringVar(&config.version, "version", "", "version of the api")
	fs.StringVar(&config.format, "format", "yaml", "format of the api content, json or yaml")

	return &ffcli.Command{
		Name:    "sync",
		FlagSet: fs,
		Exec: func(ctx context.Context, args []string) error {
			return syncExec(ctx, config)
		},
	}
}

func syncExec(ctx context.Context, config *getConfig) error {
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
			version := config.version
			if version == "" {
				var err error
				version, err = config.client.DomainSettingsDefaultGet(&swaggerhub.DomainSettingsDefaultGetParam{
					Owner:  config.owner,
					Domain: file.Name,
				})
				if err != nil {
					return fmt.Errorf("error fetching default version for %s/%s: %v", config.owner, file.Name, err)
				}
			}

			b, err := config.client.DomainGet(swaggerhub.DomainGetParam{
				Owner:       config.owner,
				Domain:      file.Name,
				Version:     version,
				ContentType: swaggerhub.ContentType{Response: "application/" + file.Format},
			})
			if err != nil {
				return fmt.Errorf("error fetching %s/%s/%s: %v", config.owner, file, version, err)
			}

			if err := handy.SaveFile(file.Path, b); err != nil {
				return fmt.Errorf("error saving file(%s): %v", file.Path, err)
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
