package confreader

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

// Функция звгрузки конфига
func LoadConf(path string) (Config, error) {
	var conf Config

	data, err := os.ReadFile(path)

	if err != nil {
		log.Errorf("Ошибка загрузки файла по пути %s : %s", path, err.Error())
		return conf, err
	}

	err = json.Unmarshal(data, &conf)
	if err != nil {
		log.Errorf("Ошибка преобразования данных файла в структуру %s", err.Error())
		return conf, err
	}

	return conf, err
}
