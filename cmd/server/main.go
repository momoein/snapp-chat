package main

import (
	"flag"
	"log"
	"os"
	"snappchat/api/http"
	app "snappchat/app/server"
	config "snappchat/config/server"
)

var configPath = flag.String("config", "server_config.json", "configuration file")

func main() {
	flag.Parse()

	if v := os.Getenv("CONFIG_PATH"); len(v) > 0 {
		*configPath = v
	}

	cfg := config.MustReadConfig(*configPath)

	appContainer := app.MustNewApp(cfg)

	log.Fatal(http.Run(appContainer))
}
