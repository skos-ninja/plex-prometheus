package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"sync"

	"github.com/jrudio/go-plex-client"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

const maxWorkers = 5

func libraryStats(p *plex.Plex) error {
	log.Println("Getting library stats")
	sections, err := p.GetLibraries()
	if err != nil {
		return err
	}

	g := errgroup.Group{}
	for _, directory := range sections.MediaContainer.Directory {
		directory := directory
		g.Go(func() error {
			return trackLibrarySize(p, directory)
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	log.Println("Finished getting library stats")

	return nil
}

func trackLibrarySize(p *plex.Plex, directory plex.Directory) error {
	results, err := p.GetLibraryContent(directory.Key, "")
	if err != nil {
		return err
	}

	// Set the size of the content in the library.
	librarySize.WithLabelValues(directory.Key, directory.Type).Set(float64(results.MediaContainer.MediaContainer.Size))

	// Set the data size of the content in the library.
	dataSize := 0.0
	switch directory.Type {
	case "movie":
		dataSize, err = getMovieLibrarySize(results)
	case "show":
		dataSize, err = getShowLibrarySize(p, results)
	}
	if err != nil {
		return err
	}
	libraryDataSize.WithLabelValues(directory.Key, directory.Type).Set(dataSize)

	return nil
}

func getShowLibrarySize(p *plex.Plex, results plex.SearchResults) (size float64, err error) {
	mu := sync.Mutex{}
	sem := semaphore.NewWeighted(maxWorkers)

	g, ctx := errgroup.WithContext(context.Background())
	for _, section := range results.MediaContainer.MediaContainer.Metadata {
		// Allow us to fetch multiple shows at once.
		err := sem.Acquire(ctx, 1)
		if err != nil {
			fmt.Printf("Acquire err = %+v\n", err)
			continue
		}

		section := section
		g.Go(func() error {
			defer sem.Release(1)
			children, err := p.GetMetadataChildren(section.RatingKey)
			if err != nil {
				return err
			}

			subG, _ := errgroup.WithContext(ctx)
			for _, child := range children.MediaContainer.Metadata {
				child := child
				subG.Go(func() error {
					seasonSize := 0.0
					episodes, err := p.GetEpisodes(child.RatingKey)
					if err != nil {
						return err
					}

					for _, episode := range episodes.MediaContainer.Metadata {
						s, err := getMediaSize(episode)
						if err != nil {
							return err
						}

						seasonSize += s
					}

					mu.Lock()
					size += seasonSize
					mu.Unlock()

					return nil
				})
			}

			if err := subG.Wait(); err != nil {
				return err
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return 0.0, err
	}

	return size, nil
}

func getMovieLibrarySize(results plex.SearchResults) (float64, error) {
	size := 0.0
	for _, meta := range results.MediaContainer.MediaContainer.Metadata {
		s, err := getMediaSize(meta)
		if err != nil {
			return 0.0, err
		}

		size += s
	}

	return size, nil
}

func getMediaSize(metadata plex.Metadata) (float64, error) {
	size := 0.0

	if len(metadata.Media) != 0 {
		for _, media := range metadata.Media {
			for _, part := range media.Part {
				size += math.Round(float64(part.Size) / 1024 / 1024)
			}
		}
	}

	return size, nil
}
