package handlers

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	botTools "github.com/cody6750/discordbot/pkg/tools"
)

//Disconnect ..
func Disconnect(s *discordgo.Session, m *discordgo.Disconnect) {
	if mediaPath != "" {
		r, err := os.Open(byeGifFile)
		if err != nil {
			log.Fatal("Could not open file")
		}
		formattedTime := botTools.CurrentTime()
		_, err = s.ChannelFileSendWithMessage("", "Mochi bot has shut down at : "+formattedTime, "Bye.gif", r)
		if err != nil {
			fmt.Printf("Error: %v", err)
		}
	}
}
