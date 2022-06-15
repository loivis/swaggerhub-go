package apis

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/loivis/swaggerhub-go"
	"github.com/loivis/swaggerhub-go/handy"
	"github.com/peterbourgon/ff/v3/ffcli"
	"golang.org/x/sync/errgroup"
)

type resolveConfig struct {
	client *swaggerhub.Client

	owner string
	paths string

	version string
	format  string
	flatten bool

	out string
}

func (c *resolveConfig) validate() error {
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

func resolve(client *swaggerhub.Client) *ffcli.Command {
	config := &resolveConfig{
		client: client,
	}

	fs := flag.NewFlagSet("swagger apis get", flag.ExitOnError)
	fs.StringVar(&config.owner, "owner", "", "owner of the apis, organization or user, case-sensitive")
	fs.StringVar(&config.paths, "paths", "", "paths of api/apis, comma separated list of files or directories")
	fs.StringVar(&config.version, "version", "", "version of the api")
	fs.StringVar(&config.format, "format", "yaml", "format of the api content, json or yaml")
	fs.BoolVar(&config.flatten, "flatten", false, "replaces all complex inline schemas with named entries in the components/schemas or definitions section")
	fs.StringVar(&config.out, "out", "resolved", "folder to save resolved APIs, with original path structure preserved")

	return &ffcli.Command{
		Name:    "resolve",
		FlagSet: fs,
		Exec: func(ctx context.Context, args []string) error {
			return resolveExec(ctx, config)
		},
	}
}

func resolveExec(ctx context.Context, config *resolveConfig) error {
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
				version, err = config.client.APISettingsDefaultGet(&swaggerhub.APISettingsDefaultGetParam{
					Owner: config.owner,
					API:   file.Name,
				})
				if err != nil {
					return fmt.Errorf("error fetching default version for %s/%s: %v", config.owner, file.Name, err)
				}
			}

			format := file.Format
			if config.format != "" {
				format = config.format
			}

			b, err := config.client.APIGet(swaggerhub.APIGetParam{
				Owner:       config.owner,
				API:         file.Name,
				Version:     version,
				ContentType: swaggerhub.ContentType{Response: "application/" + file.Format},
				Resolved:    true,
				Flatten:     config.flatten,
			})
			if err != nil {
				return fmt.Errorf("error fetching %s/%s/%s: %v", config.owner, file.Name, version, err)
			}

			path := filepath.Join(config.out, file.Path)
			path = strings.ReplaceAll(path, "."+file.Format, "."+format)
			if err := handy.SaveFile(path, b); err != nil {
				return fmt.Errorf("error saving file(%s): %v", path, err)
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
