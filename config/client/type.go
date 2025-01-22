package config

type Config struct {
	WebsocketAddr serviceAddress `json:"websocketAddr"`
}

type serviceAddress struct {
	Host string `json:"host"`
	Port uint   `json:"port"`
	Path string `json:"path"`
}
