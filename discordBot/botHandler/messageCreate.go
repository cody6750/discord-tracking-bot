package botFunctions

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

var (
	byeGifFilePath string = "mediaLibrary/Bye.gif"
	err            error
)

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
	if m.Content == "hi ms.mochi" {
		s.ChannelMessageSend(m.ChannelID, "MEOOOOOOOOW")
	}
	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}

	// If the message is "Bye" or "Goodnight" reply with GIF and message
	if m.Content == "Bye" || m.Content == "Goodnight" {
		r, err := os.Open(byeGifFilePath)
		if err != nil {
			log.Fatalf("Invalid bot parameters: %v", err)
		}
		s.ChannelFileSendWithMessage(m.ChannelID, "BYE!", "Mochi.gif", r)
	}

	if m.Content == "!Shutdown" {
		s.Close()
		log.Printf("Shutting down Mochi Both :(")
		return
	}
}
