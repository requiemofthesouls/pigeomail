package receiver

type Config struct {
	Domain            string `yaml:"domain"`
	Addr              string `yaml:"addr"`
	ReadTimeout       int    `yaml:"read_timeout_seconds"`
	WriteTimeout      int    `yaml:"write_timeout_seconds"`
	MaxMessageBytes   int    `yaml:"max_message_bytes"`
	MaxRecipients     int    `yaml:"max_recipients"`
	AllowInsecureAuth bool   `yaml:"allow_insecure_auth"`
}
