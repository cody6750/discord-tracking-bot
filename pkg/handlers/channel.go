package handlers

import "fmt"

var (
	startTracking chan struct{}
	stopTracking  chan struct{}
)

func SetChannel(channelName string, c chan struct{}) error {
	switch {
	case "startTracking" == channelName:
		startTracking = c
	case "stopTracking" == channelName:
		stopTracking = c
	default:
		return fmt.Errorf("channel %v is not a supported channel. Unable to set", channelName)
	}
	return nil
}
