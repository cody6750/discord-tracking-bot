package functions

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

//GetChannels ...
func GetChannels(s *discordgo.Session, channelType string) []*discordgo.Channel {
	var getChannels []*discordgo.Channel
	for _, guild := range s.State.Guilds {
		channels, _ := s.GuildChannels(guild.ID)
		for _, channel := range channels {
			if strings.Contains(strings.ToLower(channel.Name), channelType) {
				getChannels = append(getChannels, channel)

			}
		}
	}
	return getChannels
}

//GetDebugChannel ...
func GetDebugChannel(s *discordgo.Session) *discordgo.Channel {
	c := GetChannels(s, "debug")
	if len(c) > 0 {
		return c[0]
	}
	return nil
}