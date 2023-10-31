package utils

import (
	"os"
)

// Функция для загрузки переменной окружения по ключу
func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return ""
}
