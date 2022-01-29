package telegram

type Config struct {
	Debug   bool    `yaml:"debug"`
	Token   string  `yaml:"token"`
	Webhook Webhook `yaml:"webhook_mode"`
}

type Webhook struct {
	Enabled bool   `yaml:"enabled"`
	Port    int    `yaml:"port"`
	Cert    string `yaml:"cert"`
	Key     string `yaml:"key"`
}
