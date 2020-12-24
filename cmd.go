package main

import (
	"github.com/spf13/cobra"
)

var (
	plexAddress string
	plexToken   string
	interval    int

	cmd = &cobra.Command{
		Use:  "plex-prometheus",
		RunE: runE,
	}
)

func init() {
	cmd.PersistentFlags().StringVarP(&plexAddress, "address", "a", "http://localhost:32400", "Set the active plex server address")
	cmd.PersistentFlags().StringVarP(&plexToken, "token", "t", "", "Set the active plex token")
	cmd.PersistentFlags().IntVarP(&interval, "interval", "i", 10, "")
}
