package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/as-ifn-at/REST/config"
	"github.com/as-ifn-at/REST/internal/db/gormdbwrapper"
	"github.com/as-ifn-at/REST/internal/routes"
	"github.com/rs/zerolog"
)

func main() {

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	config := config.Load()

	dbWrapper, err := gormdbwrapper.NewDBWrapper(logger, config.DBConfigOptions, config.CacheConfig)
	if err != nil {
		panic(err.Error())
	}
	router := routes.NewRouter(config, logger, dbWrapper).SetRouters()
	listenPort := fmt.Sprintf(":%v", config.Port)
	httpServer := &http.Server{
		Addr:    listenPort,
		Handler: router,
	}

	logger.Info().Msg(fmt.Sprintf("starting server on port %v", listenPort))
	if err := httpServer.ListenAndServe(); err != nil {
		logger.Error().Msg(fmt.Sprintf("unable to start server on port %v", listenPort))
		panic(err)
	}
}
