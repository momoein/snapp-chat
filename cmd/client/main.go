package main

import (
	"flag"
	"os"
	"snappchat/api/cli"
	app "snappchat/app/client"
	config "snappchat/config/client"
)

var configPath = flag.String("config", "client_config.json", "configuration file")

func main() {
	flag.Parse()

	if v := os.Getenv("CONFIG_PATH"); len(v) > 0 {
		*configPath = v
	}
	cfg := config.MustReadConfig(*configPath)

	appContainer := app.MustNewApp(cfg)

	cli.Run(appContainer)
}
