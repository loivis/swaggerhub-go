package root

import (
	"context"
	"flag"

	"github.com/loivis/swaggerhub-go"
	"github.com/peterbourgon/ff/v3/ffcli"
)

type Config struct {
	Client *swaggerhub.Client
}

func New() (*ffcli.Command, *Config) {
	fs := flag.NewFlagSet("swagger", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "swagger",
		ShortUsage: "swagger [flags] <subcommand>",
		FlagSet:    fs,
		Exec: func(context.Context, []string) error {
			return flag.ErrHelp
		},
	}, &Config{}
}
