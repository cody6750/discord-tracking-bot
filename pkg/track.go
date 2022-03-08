package trackingbot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/cody6750/discordbot/pkg/functions"
	"github.com/cody6750/discordbot/pkg/handlers"
)

var (
	//TotalMetrics ...
	TotalMetrics          handlers.Metrics = handlers.Metrics{URL: "Total Metrics for Crawl"}
	TrackingChannelPrefix                  = "tracking_"
)

type stop struct {
	error
}

//Response ...
type Response struct {
	WebScraperResponses []handlers.ScrapeResponse `json:"WebScraperResponses"`
	Metrics             handlers.Metrics          `json:"Metrics"`
}

// TrackItemChannels ... begins the item track on the list of channels recieved. It seperates channels by type
// using the discord Guild Category Type Channel. All channels within a guild category are grouped within
// an item type. Within the configs tracking folder, the item type is used to determine the correct files to use
// per each channel.
//
//  parameters:
//
//  s *discordgo.Session : Establishes a session with discord bot.
//
//	channelsToTrack []*discordgo.Channel : List of channels to track
//
//  trackingConfigPath string : The path of the tracking file configs within the file system. These files are used
//  to define what to track on a given channel.
func (t *TrackingBot) TrackItemChannels(s *discordgo.Session, channelsToTrack []*discordgo.Channel, trackingConfigFilePath string, delay int) {
	if len(channelsToTrack) == 0 {
		return
	}

	var itemType string

	for {
		for _, channel := range channelsToTrack {
			if channel.Type == discordgo.ChannelTypeGuildCategory {
				itemType = strings.Replace(channel.Name, TrackingChannelPrefix, "", 1)
				continue
			}
			if strings.Contains(channel.Name, TrackingChannelPrefix) {
				err := t.TrackItemChannel(s, channel, trackingConfigFilePath, itemType)
				if err != nil {
					functions.LogToDiscordAndStdOut(t.logger, t.discordSession, t.discordLogChannel, t.logger.Error, err.Error())
					return
				}
				time.Sleep(1 * time.Second)
			}
			handlers.SetTotalMetrics(handlers.Metrics(TotalMetrics))
			s.ChannelMessageSend(t.discordMetricsChannel.ID, handlers.GenerateMetricsOutput(handlers.Metrics(TotalMetrics)))
			return
		}
		time.Sleep(time.Duration(int64(delay)) * time.Second)
	}
}

// TrackItemChannel ... begins the item track on the individual channel recieved. The item name to track is
// parsed from the discord channel name. The item type and item name will be used to generate a directory path
// that will resolve in the correct tracking config file.
//
//  parameters:
//
//  s *discordgo.Session : Establishes a session with discord bot.
//
//	channelsToTrack *discordgo.Channel : The channel to track
//
//  trackingConfigPath string : The path of the tracking file configs within the file system. These files are used
//  to define what to track on a given channel.
//
//  itemType string: Item type of the item to track
func (t *TrackingBot) TrackItemChannel(s *discordgo.Session, channel *discordgo.Channel, trackingConfigFilePath, itemType string) error {
	itemName := strings.Replace(channel.Name, TrackingChannelPrefix, "", 1)
	functions.LogToDiscordAndStdOut(t.logger, t.discordSession, t.discordLogChannel, t.logger.Info, fmt.Sprintf("Tracking channel: %v | item: %v | item type: %v", channel.Name, itemName, itemType))
	directoryPath := trackingConfigFilePath + itemType + "/" + itemName + "/"
	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := directoryPath + file.Name()
		err = t.trackItem(s, channel, filePath)
		if err != nil {
			return err
		}
	}
	functions.LogToDiscordAndStdOut(t.logger, t.discordSession, t.discordLogChannel, t.logger.Info, fmt.Sprintf("Successfully tracked channel: %v | item: %v | item type: %v", channel.Name, itemName, itemType))
	return nil
}

// trackItem ... makes the REST API call to the webcrawler using the tracking config file as a payload.
// The response is retrived and unmarshaled into a usuable response struct. Metrics, logs, and data are extracted
// from the response and ouputted to the corresponding discord channels.
//
//  parameters:
//
//  s *discordgo.Session : Establishes a session with discord bot.
//
//	channelsToTrack *discordgo.Channel : The channel to track
//
//  itemConfigFilePath string : The item config file path used as a payload within the REST API call to the webcrawler
func (t *TrackingBot) trackItem(s *discordgo.Session, channel *discordgo.Channel, itemConfigFilePath string) error {
	jsonFile, err := os.Open(itemConfigFilePath)
	if err != nil {
		functions.LogToDiscordAndStdOut(t.logger, t.discordSession, t.discordLogChannel, t.logger.Error, err.Error())
	}

	defer jsonFile.Close()

	payload, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		functions.LogToDiscordAndStdOut(t.logger, t.discordSession, t.discordLogChannel, t.logger.Error, err.Error())
	}

	resp, err := retry(3, time.Second*5, func() (*http.Response, error) {
		resp, err := functions.CallAPI("GET", "http://"+t.options.WebcrawlerHost+":9090/crawler/item", string(payload))
		if err != nil {
			functions.LogToDiscordAndStdOut(t.logger, t.discordSession, t.discordLogChannel, t.logger.Error, err.Error()+"Retrying....")
		}
		return resp, err
	})

	if err != nil {
		functions.LogToDiscordAndStdOut(t.logger, t.discordSession, t.discordLogChannel, t.logger.Error, err.Error()+"Retrying....")
		return err
	}

	scrapeResponse, err := t.unmarshalTrackItemResponse(resp)
	if err != nil {
		functions.LogToDiscordAndStdOut(t.logger, t.discordSession, t.discordLogChannel, t.logger.Error, err.Error())
	}

	handlers.SetCurrentMetrics(scrapeResponse.Metrics)
	s.ChannelMessageSend(t.discordMetricsChannel.ID, handlers.GenerateMetricsOutput(scrapeResponse.Metrics))
	handlers.AppendMetrics(&TotalMetrics, &scrapeResponse.Metrics)
	for _, response := range scrapeResponse.WebScraperResponses {
		for _, item := range response.ExtractedItem {
			if price, exist := item.ItemDetails["price"]; exist {
				s.ChannelMessageSend(channel.ID, item.ItemDetails["link"]+"\n```Title: "+item.ItemDetails["title"]+"\nPrice: "+price+"```")
			}
		}
	}
	return nil
}

// unmarshalTrackItemResponse takes the http.response from the REST API call to the webcrawler
// and unmarshals it into the usuauble respone struct.
func (t *TrackingBot) unmarshalTrackItemResponse(resp *http.Response) (Response, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	var res Response
	err = json.Unmarshal(body, &res)
	if err != nil {
		return Response{}, err
	}

	return res, nil
}

// retry allow users to retry a failing function given an error response.
func retry(attempts int, sleep time.Duration, fn func() (*http.Response, error)) (*http.Response, error) {
	resp, err := fn()
	if err != nil {
		if s, ok := err.(stop); ok {
			// Return the original error for later checking
			return resp, s.error
		}

		if attempts--; attempts > 0 {
			time.Sleep(sleep)
			return retry(attempts, 2*sleep, fn)
		}
		return resp, err
	}
	return resp, nil
}
