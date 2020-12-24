package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	librarySize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "plex",
		Name:      "library_item_count",
		Help:      "The total number of items in the library",
	}, []string{"library_id", "library_type"})
	libraryDataSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "plex",
		Name:      "library_data_megabytes",
		Help:      "The total data size of a plex library",
		// ConstLabels: prometheus.Labels{},
	}, []string{"library_id", "library_type"})
	activeStreams = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "plex",
		Name:      "active_streams",
		Help:      "The number of active streams",
	}, []string{"library_id", "library_type"})
)
