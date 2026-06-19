package config

import (
	"os"

	"service-controller-notebookum/internal/consul"
)

type Config struct {
	Port              string
	CorrelationHeader string
	ExtractorURL      string
	AIURL             string
	PersistenceURL    string
	UserServiceURL    string
	RedisHost         string
	RedisPort         string
	RedisPassword     string
	ConsulURL         string
}

func Load() Config {
	consulURL := env("CONSUL_URL", "http://consul:8500")
	kv := func(key, def string) string {
		return consul.KVGet(consulURL, key, def)
	}

	return Config{
		Port:              kv("port", "5000"),
		CorrelationHeader: "X-Correlation-ID",
		ExtractorURL:      kv("extractor_url", "http://extractor.universidad.localhost:5000"),
		AIURL:             kv("ai_url", "http://ai.universidad.localhost:5000"),
		PersistenceURL:    kv("persistence_url", "http://persistence-java.universidad.localhost:8080"),
		UserServiceURL:    kv("user_service_url", "http://users.universidad.localhost:5000"),
		RedisHost:         kv("redis_host", "redis"),
		RedisPort:         kv("redis_port", "6379"),
		RedisPassword:     kv("redis_password", ""),
		ConsulURL:         consulURL,
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
