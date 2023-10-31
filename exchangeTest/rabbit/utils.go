package rabbit

import (
	"fmt"
	"os"
)

// Функция для получения URL строки, для подключения к RabbitMQ
func getURLForRabbit() (string, error) {
	var envsRes map[string]string

	envs := []string{
		"RABBIT_HOST",
		"RABBIT_PORT",
		"RABBIT_USER",
		"RABBIT_PASSWORD",
		"RABBIT_VHOST",
	}

	for _, envVar := range envs {
		value, exists := os.LookupEnv(envVar)
		if !exists {
			return "", fmt.Errorf("Environment variable not loaded: %s", envVar)
		}
		envsRes[envVar] = value
	}

	result := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/%s",
		envsRes["RABBIT_USER"],
		envsRes["RABBIT_PASSWORD"],
		envsRes["RABBIT_HOST"],
		envsRes["RABBIT_PORT"],
		envsRes["RABBIT_VHOST"],
	)
	return result, nil
}
