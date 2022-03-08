package main

import (
	discordbot "github.com/cody6750/discordbot/pkg"
)

func main() {
	trackingRTX := discordbot.NewTrackingBot()
	trackingRTX.Run()
}
