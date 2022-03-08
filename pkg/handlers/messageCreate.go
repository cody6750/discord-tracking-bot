package handlers

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/cody6750/discordbot/pkg/functions"
	botTools "github.com/cody6750/discordbot/pkg/tools"
)

var (
	byeGifFile        string = "Bye.gif"
	availableCommands        = map[string]string{
		"!help":           "Displays all avaliable commands within the console channel",
		"!metrics":        "Displays the metrics for the current scrape",
		"!shutdown":       "Shuts down the bot",
		"!status":         "Displays the current status for the tracking bot",
		"!start_tracking": "Starts tracking bot using tracking configs",
		"!stop_tracking":  "Stops tracking bot",
		"!total_metrics":  "Displays the total metrics for the web crawl",
	}
)

// MessageCreate .. This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	content := strings.ToLower(m.Content)
	consoleChannel := functions.GetChannel(s, consoleChannel)
	switch {
	case content == "!help":
		helpMessage(consoleChannel, s)

	case content == "!status":
		_, err := s.ChannelMessageSend(consoleChannel.ID, "Current Status :"+currentStatus)
		if err != nil {
			logger.WithError(err).Errorf("Unable to message to channel %v", m.ChannelID)
			return
		}

	case content == "!bye" || content == "!goodnight":
		r, err := os.Open(byeGifFile)
		if err != nil {
			logger.WithError(err).Error("Could not open file")
			return
		}

		_, err = s.ChannelFileSendWithMessage(m.ChannelID, "BYE!", "Mochi.gif", r)
		if err != nil {
			logger.WithError(err).Errorf("Unable to message to channel %v", m.ChannelID)
			return
		}

	case content == "!metrics":
		if reflect.DeepEqual(currentMetrics, Metrics{}) {
			_, err := s.ChannelMessageSend(consoleChannel.ID, "No current metrics available")
			if err != nil {
				logger.WithError(err).Errorf("Unable to message to channel %v", m.ChannelID)
				return
			}
		}
		_, err := s.ChannelMessageSend(consoleChannel.ID, GenerateMetricsOutput(currentMetrics))
		if err != nil {
			logger.WithError(err).Errorf("Unable to message to channel %v", m.ChannelID)
			return
		}

	case content == "!total_metrics":
		if reflect.DeepEqual(totalMetrics, Metrics{}) {
			_, err := s.ChannelMessageSend(consoleChannel.ID, "No total metrics available")
			if err != nil {
				logger.WithError(err).Errorf("Unable to message to channel %v", m.ChannelID)
				return
			}
		}

		s.ChannelMessageSend(consoleChannel.ID, GenerateMetricsOutput(totalMetrics))

	case content == "!start_tracking":
		startTracking <- struct{}{}
		_, err := s.ChannelMessageSend(consoleChannel.ID, "Started bot tracking")
		if err != nil {
			logger.WithError(err).Errorf("Unable to message to channel %v", m.ChannelID)
			return
		}

	case content == "!stop_tracking":
		stopTracking <- struct{}{}
		_, err := s.ChannelMessageSend(consoleChannel.ID, "Stopped bot tracking")
		if err != nil {
			logger.WithError(err).Errorf("Unable to message to channel %v", m.ChannelID)
			return
		}

	case content == "!shutdown":
		ShutDownMessage(consoleChannel, s)
		s.Close()
		os.Exit(1)
	}
}

func StartUpMessage(channel *discordgo.Channel, s *discordgo.Session) {
	if mediaPath != "" {
		r, err := os.Open(mediaPath + introductionGifFile)
		if err != nil {
			logger.WithError(err).Error("Could not open file")
			return
		}

		formatedTime := botTools.CurrentTime()
		_, err = s.ChannelFileSendWithMessage(channel.ID, "Mochi Bot is now up and running at : "+formatedTime, "Introduction.gif", r)
		if err != nil {
			logger.WithError(err).Errorf("Unable to send message to channel %v", channel.ID)
		}

	} else {
		formatedTime := botTools.CurrentTime()
		_, err := s.ChannelMessageSend(channel.ID, "Mochi Bot is now up and running at : "+formatedTime)
		if err != nil {
			logger.WithError(err).Errorf("Unable to send message to channel %v", channel.ID)
		}

	}
}

func ShutDownMessage(channel *discordgo.Channel, s *discordgo.Session) {
	if mediaPath != "" {
		r, err := os.Open(mediaPath + byeGifFile)
		if err != nil {
			logger.WithError(err).Error("Could not open file")
			return
		}

		formatedTime := botTools.CurrentTime()
		_, err = s.ChannelFileSendWithMessage(channel.ID, "Mochi bot has shut down at : "+formatedTime, "Bye.gif", r)
		if err != nil {
			log.Fatal("Error sending initial message")
		}

	} else {
		formatedTime := botTools.CurrentTime()
		_, err := s.ChannelMessageSend(channel.ID, "Mochi bot has shut down at: "+formatedTime)
		if err != nil {
			log.Fatal("Error sending initial message")

		}
	}
}

func helpMessage(channel *discordgo.Channel, s *discordgo.Session) string {
	helperHeader := "```Commands available for the tracking bot:\nCommand:       Description:\n"
	body := printMap(availableCommands)
	helpMessage := helperHeader + body + "```"
	s.ChannelMessageSend(channel.ID, helpMessage)
	return helpMessage
}

func printMap(m map[string]string) string {
	var maxLenKey int
	for k, _ := range m {
		if len(k) > maxLenKey {
			maxLenKey = len(k)
		}
	}

	var mapLine string
	for k, v := range m {
		line := fmt.Sprintln(k + " " + strings.Repeat(" ", maxLenKey-len(k)) + v)
		mapLine += line
	}
	return mapLine
}
