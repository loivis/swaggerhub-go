package domains

import (
	"context"
	"flag"

	"github.com/loivis/swaggerhub-go"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func New(client *swaggerhub.Client) *ffcli.Command {
	fs := flag.NewFlagSet("swagger domains", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "domains",
		ShortUsage: "domains [flags] <subcommand>",
		FlagSet:    fs,
		Subcommands: []*ffcli.Command{
			sync(client),
		},
		Exec: func(ctx context.Context, args []string) error {
			return flag.ErrHelp
		},
	}
}
