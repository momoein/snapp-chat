package app

import (
	config "snappchat/config/server"
	"snappchat/internal/server"
)

type App interface {
	Service() server.Service
	Config() config.Config
}
