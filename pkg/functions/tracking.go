package functions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	itemName                string
	channelName             string
	channelID               string
	debugChannel            *discordgo.Channel
	introductionGifFilePath string = "media/Introduction.gif"
	guildCategory                  = map[string]string{
		"Tracking Graphics Cards": "graphic_cards",
		"Tracking Consoles":       "consoles",
	}
)

// type FilterConfiguration struct {
// 	Contains      string
// 	IsEqualTo     interface{}
// 	IsLessThan    int
// 	IsGreaterThan int
// 	IsNot         interface{}
// }

type Response []ScrapeResponse

//ScrapeResponse ...
type ScrapeResponse struct {
	RootURL       string `json:"RootURL"`
	ExtractedItem []Item `json:"ExtractedItem"`
	ExtractedURLs []URL  `json:"ExtractedURLs"`
}

//Item ...
type Item struct {
	ItemName    string            `json:"ItemName"`
	URL         URL               `json:"URL"`
	TimeQueried string            `json:"TimeQueried"`
	DateQueried string            `json:"DateQueried"`
	ItemDetails map[string]string `json:"ItemDetails"`
	// ItemFilterConfiguration FilterConfiguration `json:"FilterConfiguration"`
}

//URL ...
type URL struct {
	RootURL      string `json:"RootURL"`
	ParentURL    string `json:"ParentURL"`
	CurrentURL   string `json:"CurrentURL"`
	CurrentDepth int    `json:"CurrentDepth"`
	MaxDepth     int    `json:"MaxDepth"`
}

func getPayload() {

}

//StartTracking ...
func StartTracking(s *discordgo.Session, channels []*discordgo.Channel) {
	if len(channels) == 0 {
		return
	}

	debugChannel = GetDebugChannel(s)
	if debugChannel == nil {
		log.Fatal("Debug channel must exist")
	}
	var itemType string
	for {
		log.Print("Tracking")
		for _, channel := range channels {
			if channel.Type == discordgo.ChannelTypeGuildCategory {
				itemType = guildCategory[channel.Name]
				continue
			}
			if strings.Contains(channel.Name, "tracking") {
				log.Printf("Beginning to track channel %v, item type %v", channel.Name, itemType)
				track(s, channel, strings.Replace(channel.Name, "tracking_", "", 1), itemType)
			}
		}
		time.Sleep(6 * time.Hour)
	}
}

func track(s *discordgo.Session, channel *discordgo.Channel, itemName, itemType string) {
	directoryPath := "pkg/configs/tracking/" + itemType + "/" + itemName + "/"
	if !strings.Contains(itemName, "ti") {
		return
	}
	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		filePath := directoryPath + file.Name()

		// Open our jsonFile
		jsonFile, err := os.Open(filePath)

		// if we os.Open returns an error then handle it
		if err != nil {
			s.ChannelMessageSend(debugChannel.ID, err.Error())
		}
		// defer the closing of our jsonFile so that we can parse it later on
		defer jsonFile.Close()

		payload, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			fmt.Println(err)
		}
		log.Print("Calling payload")
		resp, err := CallAPI(string(payload))
		if err != nil {
			log.Print(err.Error())
			s.ChannelMessageSend(debugChannel.ID, err.Error())
		}
		scrapeResponse := unmarshalTrackItemResponse(resp)
		for _, response := range scrapeResponse {
			for _, item := range response.ExtractedItem {
				if price, exist := item.ItemDetails["price"]; exist {
					if err != nil {
						log.Fatal("Unable to convert string value to int")
					}
					s.ChannelMessageSend(channel.ID, item.ItemDetails["link"]+"\n```Title:"+item.ItemDetails["title"]+"\nPrice $"+price+"```")
				}
			}
		}
		s.ChannelMessageSend(channel.ID, "Succesfully tracked")
		time.Sleep(2 * time.Second)
	}
}

func unmarshalTrackItemResponse(resp *http.Response) []ScrapeResponse {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	var res Response
	json.Unmarshal(body, &res)
	return res
}
