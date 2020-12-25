package main

import (
	"log"

	"github.com/jrudio/go-plex-client"
)

func streamStats(p *plex.Plex) {
	log.Println("Registering for plex events")
	events := plex.NewNotificationEvents()

	events.OnPlaying(func(n plex.NotificationContainer) {
		state := map[string]map[string]float64{}
		sessions, err := p.GetSessions()
		if err != nil {
			panic(err)
		}

		for _, session := range sessions.MediaContainer.Metadata {
			m := state[session.Player.State]
			if m == nil {
				state[session.Player.State] = map[string]float64{}
			}
			state[session.Player.State][session.LibrarySectionID]++
		}

		// Remove all the old vectors in case any state library pair is missing
		activeStreams.Reset()
		for playerState, v := range state {
			for libraryID, d := range v {
				activeStreams.WithLabelValues(playerState, libraryID).Set(d)
			}
		}
	})

	registerEvents(p, events)
	log.Println("Registered for plex events")
}
