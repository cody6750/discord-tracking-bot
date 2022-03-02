package bothandler

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	botTools "github.com/cody6750/discordbot/pkg/tools"
)

var (
	channelID               string
	introductionGifFilePath string = "media/Introduction.gif"
	guildCategory                  = map[string]string{
		"Tracking Graphics Cards": "graphic_cards",
		"Tracking Consoles":       "consoles",
	}
)

//Ready ..
func Ready(s *discordgo.Session, m *discordgo.Ready) {
	for _, guild := range s.State.Guilds {

		channels, _ := s.GuildChannels(guild.ID)

		for _, channel := range channels {
			if channel.Type != discordgo.ChannelTypeGuildText {
				continue
			}
			r, err := os.Open(introductionGifFilePath)

			if err != nil {
				log.Fatalf("Error opening %v", introductionGifFilePath)
			}
			if channel.Name == "test" {
				formatedTime := botTools.CurrentTime()
				_, err = s.ChannelFileSendWithMessage(channel.ID, "Mochi Bot is now up and running at : "+formatedTime, "Introduction.gif", r)
				if err != nil {
					log.Fatal("Error sending initial message")
				}
			}
		}
	}
}
