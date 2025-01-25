package app

import (
	"fmt"
	"net/url"
	config "snappchat/config/client"
	"snappchat/internal/client"
)

type app struct {
	service client.Service
	cfg     config.Config
}

func (a *app) Service() client.Service {
	return a.service
}

func (a *app) Config() config.Config {
	return a.cfg
}

func (a *app) setService() error {
	if a.cfg.WebsocketAddr.IsEmpty() {
		return fmt.Errorf("WebSocket address is not configured")
	}
	wsURL := url.URL{
		Scheme: "ws",
		Host:   fmt.Sprintf("%s:%d", a.cfg.WebsocketAddr.Host, a.cfg.WebsocketAddr.Port),
		Path:   a.cfg.WebsocketAddr.Path,
	}

	if a.cfg.HttpAddr.IsEmpty() {
		return fmt.Errorf("http addr is not configured")
	}
	httpURL := url.URL{
		Scheme: "http",
		Host: fmt.Sprintf("%s:%d", a.cfg.HttpAddr.Host, a.cfg.HttpAddr.Port),
		Path: a.cfg.HttpAddr.Path,
	}

	a.service = client.NewService(wsURL.String(), httpURL.String())
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
