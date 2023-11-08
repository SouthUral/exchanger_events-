package main

import (
	"fmt"

	confReader "github.com/SouthUral/exchangeTest/confreader"
)

func main() {
	conf, err := confReader.LoadConfRabbit("./config/example.json")
	if err != nil {
		return
	}

	fmt.Println(conf)
}
