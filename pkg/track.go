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

//TrackItemChannels ...
func (t *TrackingBot) TrackItemChannels(s *discordgo.Session, channels []*discordgo.Channel, trackingConfigFilePath string, delay int) {
	if len(channels) == 0 {
		return
	}

	var itemType string

	for {
		for _, channel := range channels {
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

//TrackItemChannel ...
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
				s.ChannelMessageSend(channel.ID, item.ItemDetails["link"]+"\n```Title:"+item.ItemDetails["title"]+"\nPrice $"+price+"```")
			}
		}
	}
	return nil
}

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
