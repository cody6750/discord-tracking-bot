package handlers

import (
	"os"

	"github.com/bwmarrin/discordgo"
	botTools "github.com/cody6750/discordbot/pkg/tools"
)

// Disconnect serves as the disconnect handler. When discord is disconnecting, this function gets executed.
// A shutdown message will be called within the console channel.
func Disconnect(s *discordgo.Session, m *discordgo.Disconnect) {
	if mediaPath != "" {
		r, err := os.Open(byeGifFile)
		if err != nil {
			logger.Error("Could not open file")
			return
		}
		formattedTime := botTools.CurrentTime()
		_, err = s.ChannelFileSendWithMessage(consoleChannel, "Mochi bot has shut down at : "+formattedTime, "Bye.gif", r)
		if err != nil {
			logger.WithError(err).Errorf("unable to send message to channel: %v", consoleChannel)
			return
		}
	}
}
