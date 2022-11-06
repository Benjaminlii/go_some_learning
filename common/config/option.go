package config

type RocketMQConfig struct {
	NameServers []string                     `yaml:"name_servers"`
	Topic       string                       `yaml:"topic"`
	TagMap      map[string]RocketMQTagOption `yaml:"tag_map"`
}
type RocketMQTagOption struct {
	Tags []string `yaml:"tags"`
}

type RedisOption struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
}
