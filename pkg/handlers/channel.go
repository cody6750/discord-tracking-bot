package handlers

import "fmt"

var (
	//startTracking used to execute the start tracking function within the tracking bot from the handlers package
	startTracking chan struct{}

	//stopTracking used to execute the stop tracking function within the tracking bot from the handlers package
	stopTracking chan struct{}
)

//SetChannel used to pass in the channels from the trackingbot package to avoid circular dependency
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
