package config

type Config struct {
	WebsocketAddr serviceAddress `json:"websocketAddr"`
	HttpAddr      serviceAddress `json:"httpAddr"`
}

type serviceAddress struct {
	Host string `json:"host"`
	Port uint   `json:"port"`
	Path string `json:"path"`
}

func (s *serviceAddress) IsEmpty() bool {
	return s.Host == "" || s.Port == 0 || s.Path == ""
}
