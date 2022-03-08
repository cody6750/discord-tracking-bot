package handlers

import (
	"github.com/bwmarrin/discordgo"
)

var (
	consoleChannel      string = "bot_console"
	metricsChannel      string = "metrics"
	logChannel          string = "log"
	introductionGifFile string = "Introduction.gif"
)

//Ready ..
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
