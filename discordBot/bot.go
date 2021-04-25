package discordbot

// IF broken import, use go get github.com/bwmarrin/discordgo
// DOCCUMENTATION REFFERENCE: https://github.com/bwmarrin/discordgo/blob/master/examples/pingpong/main.go
import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	bothandler "github.com/cody6750/discordbot/discordBot/botHandler"
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
	discord.AddHandler(bothandler.Ready)
	discord.AddHandler(bothandler.MessageCreate)
	discord.AddHandler(bothandler.Disconnect)
	discord.AddHandler(bothandler.SlashCommands)

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
