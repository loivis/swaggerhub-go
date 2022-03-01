package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/loivis/swaggerhub-go"
	"github.com/loivis/swaggerhub-go/cmd/swagger/apis"
	"github.com/loivis/swaggerhub-go/cmd/swagger/root"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func main() {
	client, err := swaggerClient()
	if err != nil {
		log.Fatalf("error creating swaggerhub api client: %v", err)
	}

	rootCommand, rootConfig := root.New()

	rootCommand.Subcommands = []*ffcli.Command{
		apis.New(client),
	}

	if err := rootCommand.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error during Parse: %v\n", err)
		os.Exit(1)
	}

	rootConfig.Client = client

	if err := rootCommand.Run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func swaggerClient() (*swaggerhub.Client, error) {
	baseURL := os.Getenv("SWAGGERHUB_BASE_URL")
	apiKey := os.Getenv("SWAGGERHUB_API_KEY")

	errs := make([]string, 0, 2)

	if baseURL == "" {
		errs = append(errs, "SWAGGERHUB_BASE_URL")
	}

	if apiKey == "" {
		errs = append(errs, "SWAGGERHUB_API_KEY")
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("error reading environment variables: %s", strings.Join(errs, ", "))
	}

	return swaggerhub.New(baseURL, apiKey), nil
}
