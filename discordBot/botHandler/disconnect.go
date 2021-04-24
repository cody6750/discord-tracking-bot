package bothandler

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	botTools "github.com/cody6750/RTXBot/discordBot/botTools"
)

//Disconnect ..
func Disconnect(s *discordgo.Session, m *discordgo.Disconnect) {
	r, err := os.Open(byeGifFilePath)
	if err != nil {
		log.Fatal("Could not open file")
	}
	formattedTime := botTools.CurrentTime()
	_, err = s.ChannelFileSendWithMessage("", "Mochi bot has shut down at : "+formattedTime, "Bye.gif", r)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}
