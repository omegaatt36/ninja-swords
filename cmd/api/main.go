package main

import (
	"context"

	"github.com/urfave/cli/v2"
	"ninja-swords/app"
	"ninja-swords/app/api"
	"ninja-swords/health"
)

// Main starts process in cli.
func Main(ctx context.Context, c *cli.Context) {
	go health.StartServer()

	server := api.Server{}
	server.Start(ctx, c.String("listen-addr"))
}

func main() {
	app := app.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "listen-addr",
				Value: ":8787",
			},
		},
		Main: Main,
	}

	app.Run()
}
