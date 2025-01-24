package app

import (
	config "snappchat/config/client"
	"snappchat/internal/client"
)

type App interface {
	Service() client.Service
	Config() config.Config
}
