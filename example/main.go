package main

import (
	"github.com/obnahsgnaw/application"
	"time"
)

func main() {
	app := application.New(
		application.NewCluster("dev", "Dev"),
		"demo",
		application.Debug(func() bool { return true }),
	)
	defer app.Release()

	app.With(application.EtcdRegister([]string{"127.0.0.1:2379"}, 5*time.Second))

	app.Run(func(err error) {
		panic(err)
	})
	app.Wait()
}
