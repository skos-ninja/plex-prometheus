package main

import (
	"errors"
	"os"

	"github.com/jrudio/go-plex-client"
)

func newConnection() (*plex.Plex, error) {
	if plexToken == "" {
		return nil, errors.New("missing plex token")
	}

	connection, err := plex.New(plexAddress, plexToken)
	if err != nil {
		return nil, err
	}
	connection.Headers.Platform = "plex-prometheus"
	connection.Headers.Product = "Plex Prometheus Metrics"

	valid, err := connection.Test()
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("failed to connect to plex")
	}

	return connection, nil
}

func registerEvents(p *plex.Plex, e *plex.NotificationEvents) {
	ctrlC := make(chan os.Signal, 1)
	p.SubscribeToNotifications(e, ctrlC, func(err error) {
		panic(err)
	})
}
