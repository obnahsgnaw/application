package main

import (
	"github.com/obnahsgnaw/application"
	"github.com/obnahsgnaw/application/pkg/debug"
	"log"
	"time"
)

func main() {
	app := application.New(
		"demo",
		"demo",
		application.Debugger(debug.New(true)),
	)
	app.With(application.EtcdRegister([]string{"127.0.0.1:2379"}, 5*time.Second))

	app.Run(func(err error) {
		panic(err)
	})
	log.Println("Exited")
}
