package telegram

type (
	Config struct {
		Debug   bool          `mapstructure:"debug"`
		Token   string        `mapstructure:"token"`
		Webhook WebhookConfig `mapstructure:"webhook"`
	}
	WebhookConfig struct {
		Enabled bool   `mapstructure:"enabled"`
		Port    uint32 `mapstructure:"port"`
		Cert    string `mapstructure:"cert"`
		Key     string `mapstructure:"key"`
		Domain  string `mapstructure:"domain"`
	}
)
