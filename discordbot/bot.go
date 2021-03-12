package discordbot

// IF broken import, use go get github.com/bwmarrin/discordgo
// DOCCUMENTATION REFFERENCE: https://github.com/bwmarrin/discordgo/blob/master/examples/pingpong/main.go
import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	// Need to store token in secrets manager or something
	discordToken string = "ODE5NDgyNjU1NjczODEwOTU0.YEnQsg.Pmu9ppRwPpvJl6ebPS2ryQCyrnc"
)

//BotInit initializes the discord bot
func BotInit() {

	fmt.Println("Discord bot initializing......")

	// Create a new Discord session using the provided bot token.
	discord, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Println("Error has occured")
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	discord.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	discord.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
