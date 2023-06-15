package main

import (
	"github.com/obnahsgnaw/application"
	"github.com/obnahsgnaw/application/pkg/debug"
	"log"
)

func main() {
	app := application.New(
		"demo",
		"demo",
		application.Debugger(debug.New(true)),
	)

	app.Run(func(err error) {
		panic(err)
	})
	log.Println("Exited")
}
