package config

type Config struct {
	BrokerURL string
	Topic     string
	Group     string
}

func NewConfig() *Config {
	return &Config{
		BrokerURL: "nats://127.0.0.1:4222",
		Topic:     "MyJSTopic",
		Group:     "MyQueueGroup",
	}
}
