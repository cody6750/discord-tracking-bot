package main

import (
	discordbot "github.com/cody6750/discordbot/pkg"
)

var (
	globalvar string
)

func main() {
	trackingRTX := discordbot.NewTrackingBot()
	trackingRTX.Run()
}
