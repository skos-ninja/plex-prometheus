package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}

func runE(cmd *cobra.Command, args []string) error {
	plex, err := newConnection()
	if err != nil {
		return err
	}

	go func() {
		for {
			if err := libraryStats(plex); err != nil {
				panic(err)
			}

			time.Sleep(time.Second * time.Duration(interval))
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(":8080", nil)
}
