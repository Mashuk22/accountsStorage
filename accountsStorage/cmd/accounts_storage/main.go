package main

import (
	"account_storage/internal/app/apiserver"
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

func main() {

	var (
		configPath string
		logger     *logrus.Logger
	)

	flag.StringVar(&configPath, "config-path", "./cmd/accounts_storage/configs/apiserver.toml", "path to the config file")
	flag.Parse()

	logger = logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	ctx, cancel := context.WithCancel(context.Background())
	go handleSignals(cancel)

	config := apiserver.NewConfig()

	_, err := toml.DecodeFile(configPath, &config)
	if err != nil {
		logger.Fatal(err)
	}

	server, err := apiserver.NewServer(logger, ctx, config)
	if err != nil {
		logger.Fatal(err)
	}

	err = server.Start()
	if err != nil {
		logger.Fatal(err)
	}
}

func handleSignals(cancel context.CancelFunc) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	<-signalChannel
	cancel()
}
