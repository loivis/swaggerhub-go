package apis

import (
	"context"
	"flag"

	"github.com/loivis/swaggerhub-go"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func New(client *swaggerhub.Client) *ffcli.Command {
	fs := flag.NewFlagSet("swagger apis", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "apis",
		ShortUsage: "apis [flags] <subcommand>",
		FlagSet:    fs,
		Subcommands: []*ffcli.Command{
			sync(client),
			resolve(client),
			update(client),
			setDefault(client),
			publish(client),
		},
		Exec: func(ctx context.Context, args []string) error {
			return flag.ErrHelp
		},
	}
}
