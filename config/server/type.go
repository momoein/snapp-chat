package config

type Config struct {
	Server ServerConfig `json:"server"`
	Nats   NatsConfig   `json:"nats"`
}

type ServerConfig struct {
	HttpPort          uint   `json:"httpPort"`
	Secret            string `json:"secret"`
	AuthExpMinute     uint   `json:"authExpMin"`
	AuthRefreshMinute uint   `json:"authExpRefreshMin"`
}

type NatsConfig struct {
	Host string `json:"host"`
	Port uint   `json:"port"`
}

type ClientConfig struct {
	Host string `json:"host"`
	Port uint   `json:"port"`
	Path string `json:"path"`
}
