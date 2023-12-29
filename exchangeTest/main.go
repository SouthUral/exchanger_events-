package main

import (
	"os"
	"time"

	router "github.com/SouthUral/exchangeTest/router"
	ut "github.com/SouthUral/exchangeTest/router/utilsfortest"

	log "github.com/sirupsen/logrus"
)

func init() {
	// логи в формате JSON, по умолчанию формат ASCII
	log.SetFormatter(&log.JSONFormatter{})

	// логи идут на стандартный вывод, их можно перенаправить в файл
	log.SetOutput(os.Stdout)

	// установка уровня логирования
	log.SetLevel(log.DebugLevel)
}

func main() {
	conf, err := ut.LoadConf("./router/testdata/fixt_1.json")
	if err != nil {
		return
	}
	log.Debug(conf)

	eventCh, subscrCh, cancelRouter := router.InitRouter()
	ut.StartPublishers(conf.Publishers, eventCh, 50)
	ut.SubscribersWork(conf.Consumers, subscrCh, 50)

	time.Sleep(5 * time.Minute)
	cancelRouter()
	time.Sleep(30 * time.Second)
}
