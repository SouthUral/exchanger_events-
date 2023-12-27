package utilsfortest

import (
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// Функция звгрузки конфига
func LoadConf(path string) (Config, error) {
	var conf Config

	data, err := os.ReadFile(path)

	if err != nil {
		err = fmt.Errorf("Ошибка загрузки файла по пути %s : %s", path, err.Error())
		log.Error(err)
		return conf, err
	}

	err = json.Unmarshal(data, &conf)
	if err != nil {
		err = fmt.Errorf("Ошибка преобразования данных файла в структуру %s", err.Error())
		log.Error(err)
		return conf, err
	}

	return conf, err
}
