package handlers

import (
	"github.com/bwmarrin/discordgo"
)

var (
	//consoleChannel name of the console channel in Discord.
	consoleChannel string = "bot_console"

	//consoleChannel name of the metrics channel in Discord.
	metricsChannel string = "metrics"

	//consoleChannel name of the log channel in Discord.
	logChannel string = "log"

	//introductionGifFile file name of the introduction gif.
	introductionGifFile string = "Introduction.gif"

	//byeGifFile file name of the bye gif.
	byeGifFile string = "Bye.gif"
)

// Ready serves as the ready handler. When discord is in a ready state, this function gets executed.
// A startup message function will be called within the console channel.
func Ready(s *discordgo.Session, m *discordgo.Ready) {
	for _, guild := range s.State.Guilds {
		channels, _ := s.GuildChannels(guild.ID)
		for _, channel := range channels {
			if channel.Type != discordgo.ChannelTypeGuildText {
				continue
			}
			if channel.Name == consoleChannel {
				StartUpMessage(channel, s)
			}
		}
	}
}
