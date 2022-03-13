package functions

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

//GetChannels retrives channels given the channel type within the Discord guild.
func GetChannels(s *discordgo.Session, channelType string) []*discordgo.Channel {
	var getChannels []*discordgo.Channel
	for _, guild := range s.State.Guilds {
		channels, err := s.GuildChannels(guild.ID)
		if err != nil {
			log.Fatalf(err.Error())
		}
		for _, channel := range channels {
			if channel == nil {
				continue
			}
			if strings.Contains(strings.ToLower(channel.Name), channelType) {
				getChannels = append(getChannels, channel)
			}
		}
	}
	return getChannels
}

//GetChannel retrives channel given the channel name within the Discord guild.
func GetChannel(s *discordgo.Session, channelName string) *discordgo.Channel {
	c := GetChannels(s, channelName)
	if len(c) > 0 {
		return c[0]
	}
	return nil
}
