package main

import (
	"github.com/obnahsgnaw/application"
	"github.com/obnahsgnaw/application/pkg/logging/logger"
	"time"
)

func main() {
	app := application.New(
		application.NewCluster("dev", "Dev"),
		"demo",
		application.Debug(func() bool { return true }),
	)
	defer app.Release()

	app.With(application.Logger(&logger.Config{
		Dir:        "/Users/wangshanbo/Documents/Data/projects/application/out",
		MaxSize:    100,
		MaxBackup:  5,
		MaxAge:     10,
		Level:      "debug",
		TraceLevel: "error",
	}))

	app.With(application.EtcdRegister([]string{"127.0.0.1:2379"}, 5*time.Second))

	app.Run(func(err error) {
		panic(err)
	})
	app.Wait()
}
