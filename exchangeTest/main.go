package main

import (
	"os"
	// "time"

	// router "github.com/SouthUral/exchangeTest/router"

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
	// conf, err := confReader.LoadConf("./config/example.json")
	// if err != nil {
	// 	return
	// }

	// eventCh, subscrCh, cancelRouter := router.InitRouter()
	// cancelPub := pub.StartPublishers(conf.Publishers, eventCh)
	// cancelSub := sub.StartSubscribers(conf.Consumers, subscrCh)

	// time.Sleep(5 * time.Minute)
	// cancelSub()
	// cancelPub()
	// cancelRouter()
	// time.Sleep(30 * time.Second)
}
