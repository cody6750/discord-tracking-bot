package handlers

import (
	"os"

	"github.com/bwmarrin/discordgo"
	botTools "github.com/cody6750/discordbot/pkg/tools"
)

//Disconnect ..
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
