package functions

import "github.com/bwmarrin/discordgo"

func Ready(s *discordgo.Session, m *discordgo.Ready) {
	s.ChannelMessage("819444650972413992", "hello")
}
