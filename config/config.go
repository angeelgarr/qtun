package config

type Config struct {
	From       string
	To         string
	ClientKey  string
	ClientPem  string
	ServerKey  string
	ServerPem  string
	Timeout    int
	ServerMode bool
	UDPMode    bool
}
