package config

type AppConfig struct {
	RocketMQ *RocketMQConfig `yaml:"rocket_mq"`
	Redis    *RedisOption    `yaml:"redis"`
}
