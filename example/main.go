package main

import (
	"github.com/obnahsgnaw/application"
	"github.com/obnahsgnaw/application/pkg/logging/logger"
	"github.com/obnahsgnaw/application/service/regCenter"
	"time"
)

func main() {
	app := application.New("demo")
	defer app.Release()

	app.With(application.CusCluster(application.NewCluster("dev", "Dev")))
	app.With(application.Debug(func() bool { return true }))
	app.With(application.Logger(&logger.Config{
		Dir:        "/Users/wangshanbo/Documents/Data/projects/application/out",
		MaxSize:    100,
		MaxBackup:  5,
		MaxAge:     10,
		Level:      "debug",
		TraceLevel: "error",
	}))
	r, _ := regCenter.NewEtcdRegister([]string{"127.0.0.1:2379"}, 5*time.Second)
	app.With(application.Register(r, 5))

	app.Run(func(err error) {
		panic(err)
	})
	app.Wait()
}
