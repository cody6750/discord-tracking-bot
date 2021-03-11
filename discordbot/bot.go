package discordbot

// IF broken import, use go get github.com/bwmarrin/discordgo
import "github.com/bwmarrin/discordgo"

func bot() {
	botInit()
}

func botInit() {
	discord, err := discordgo.New("Bot " + "authentication token")
}
