package app

import (
	"fmt"
	config "snappchat/config/server"
	"snappchat/internal/server"

	"github.com/nats-io/nats.go"
)

type app struct {
	service server.Service
	cfg     config.Config
}

func (a *app) Service() server.Service {
	return a.service
}

func (a *app) Config() config.Config {
	return a.cfg
}

func (a *app) setService() error {
	natsServer := fmt.Sprintf("%s:%d", a.cfg.Nats.Host, a.cfg.Nats.Port)
	nc, err := nats.Connect(natsServer)
	if err != nil {
		return fmt.Errorf("error on connecting to nats: %v", err)
	}

	a.service = server.NewService(nc)
	return nil
}

func NewApp(cfg config.Config) (App, error) {
	a := &app{
		cfg: cfg,
	}

	if err := a.setService(); err != nil {
		return nil, err
	}

	return a, nil
}

func MustNewApp(cfg config.Config) App {
	a, err := NewApp(cfg)
	if err != nil {
		panic(err)
	}
	return a
}
