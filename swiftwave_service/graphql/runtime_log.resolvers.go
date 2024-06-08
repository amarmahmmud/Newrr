package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.48

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"time"

	dbmodel "github.com/swiftwave-org/swiftwave/swiftwave_service/core"
	"github.com/swiftwave-org/swiftwave/swiftwave_service/graphql/model"
)

// FetchRuntimeLog is the resolver for the fetchRuntimeLog field.
func (r *subscriptionResolver) FetchRuntimeLog(ctx context.Context, applicationID string, timeframe model.RuntimeLogTimeframe) (<-chan *model.RuntimeLog, error) {
	// fetch application
	var application dbmodel.Application
	err := application.FindById(ctx, r.ServiceManager.DbClient, applicationID)
	if err != nil {
		return nil, err
	}
	// fetch docker manager
	dockerManager, err := FetchDockerManager(ctx, &r.ServiceManager.DbClient)
	if err != nil {
		return nil, err
	}
	// fetch runtime logs
	var sinceMinutes int
	switch timeframe {
	case model.RuntimeLogTimeframeLive:
		sinceMinutes = 1
	case model.RuntimeLogTimeframeLast1Hour:
		sinceMinutes = 60
	case model.RuntimeLogTimeframeLast3Hours:
		sinceMinutes = 180
	case model.RuntimeLogTimeframeLast6Hours:
		sinceMinutes = 360
	case model.RuntimeLogTimeframeLast12Hours:
		sinceMinutes = 720
	case model.RuntimeLogTimeframeLast24Hours:
		sinceMinutes = 1440
	case model.RuntimeLogTimeframeLifetime:
		sinceMinutes = 0
	}
	logsReader, err := dockerManager.LogsService(application.Name, sinceMinutes)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(logsReader)
	// create a channel
	var channel = make(chan *model.RuntimeLog, 200)
	// start a goroutine
	go func() {
		defer close(channel)
		// defer handle panic
		defer func() {
			if r := recover(); r != nil {
				return
			}
		}()
		// close logs reader
		defer func(logsReader io.ReadCloser) {
			err := logsReader.Close()
			if err != nil {
				fmt.Println("error while closing logs reader")
			}
		}(logsReader)
		// iterate over logs
		for scanner.Scan() {
			// Specific format for raw-stream logs
			// docs : https://docs.docker.com/engine/api/v1.42/#tag/Container/operation/ContainerAttach
			logTextBytes := scanner.Bytes()
			if len(logTextBytes) > 8 {
				logTextBytes = logTextBytes[8:]
			}
			// add new line
			logTextBytes = append(logTextBytes, []byte("\n")...)

			select {
			case <-ctx.Done():
				return
			case channel <- &model.RuntimeLog{
				Content:   string(logTextBytes),
				CreatedAt: time.Now(),
			}: // do nothing
			}
		}
	}()

	return channel, nil
}
