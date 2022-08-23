package main

import (
	v1 "github.com/XBozorg/bookstore/adapter/delivery/http/v1"
	"github.com/XBozorg/bookstore/adapter/repository"
	"github.com/XBozorg/bookstore/config"
	"github.com/XBozorg/bookstore/log"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	log.DefLoggers(&log.I, &log.E, &log.H) // initialize loggers

	log.I.Infoln("Starting...")

	err := config.Conf.Read() // read config file
	if err != nil {
		log.E.Panic(err)
	}

	log.I.Infoln("Config file Loaded")

	err = repository.Repo.Connect(&config.Conf) // connect repository to databases
	if err != nil {
		log.E.Panic(err)
	}

	log.I.Infoln("Repository Connected")

}

func main() {

	defer repository.Repo.Close()

	e := v1.Routing(repository.Repo)

	e.Use(middleware.Recover())

	e.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Format:           config.Conf.GetEchoConfig().LoggerFormat,
			CustomTimeFormat: log.TimeFMT,
			Output:           log.H.Out,
		}))

	e.Logger.Fatal(e.Start(config.Conf.GetEchoConfig().HttpPort))
}
