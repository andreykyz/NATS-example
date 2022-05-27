package config

type Config struct {
	BrokerURL      string
	Stream         string
	SubjectFilter  string
	Consumer       string
	ConsumerFilter string
}

func NewConfig() *Config {
	return &Config{
		BrokerURL:      "nats://127.0.0.1:4222",
		Stream:         "MyJStream",
		SubjectFilter:  "MyJSTopic",
		Consumer:       "MyQueueConsumer",
		ConsumerFilter: "MyFilter",
	}
}
