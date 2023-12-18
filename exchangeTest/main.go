package main

import (
	"os"
	"time"

	confReader "github.com/SouthUral/exchangeTest/confreader"
	pub "github.com/SouthUral/exchangeTest/publisher"
	router "github.com/SouthUral/exchangeTest/router"
	sub "github.com/SouthUral/exchangeTest/subscriber"

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
	conf, err := confReader.LoadConfRabbit("./config/example.json")
	if err != nil {
		return
	}

	eventCh, subscrCh, cancelRouter := router.InitRouter()
	cancelPub := pub.StartPublishers(conf.Publishers, eventCh)
	cancelSub := sub.StartSubscribers(conf.Consumers, subscrCh)

	time.Sleep(2 * time.Minute)
	cancelSub()
	cancelPub()
	cancelRouter()
}
