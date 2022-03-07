package handlers

import "strconv"

// Metrics ...
type Metrics struct {
	URL                 string
	DuplicatedUrlsFound int `json:"DuplicatedUrlsFound"`
	UrlsFound           int `json:"UrlsFound"`
	UrlsVisited         int `json:"UrlsVisited"`
	ItemsFound          int `json:"ItemsFound"`
}

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
}

//URL ...
type URL struct {
	RootURL      string `json:"RootURL"`
	ParentURL    string `json:"ParentURL"`
	CurrentURL   string `json:"CurrentURL"`
	CurrentDepth int    `json:"CurrentDepth"`
	MaxDepth     int    `json:"MaxDepth"`
}

var (
	currentMetrics Metrics
	totalMetrics   Metrics
)

func SetCurrentMetrics(m Metrics) Metrics {
	currentMetrics = m
	return currentMetrics
}

func SetTotalMetrics(m Metrics) Metrics {
	totalMetrics = m
	return totalMetrics
}

func AppendMetrics(m, a *Metrics) {
	m.DuplicatedUrlsFound += a.DuplicatedUrlsFound
	m.ItemsFound += a.ItemsFound
	m.UrlsFound += a.UrlsFound
	m.UrlsVisited += a.UrlsVisited
}

func GenerateMetricsOutput(m Metrics) string {
	Header := "Scraping Metrics :" + m.URL
	duplicateUrlsFound := "Duplicate Urls Found: " + strconv.Itoa(m.DuplicatedUrlsFound)
	urlsFound := "Urls Found: " + strconv.Itoa(m.UrlsFound)
	itemsFound := "Items Found: " + strconv.Itoa(m.ItemsFound)
	urlsVisited := "Urls visited Found: " + strconv.Itoa(m.ItemsFound)
	output := "\n```" + Header + "\n" + duplicateUrlsFound + "\n" + urlsFound + "\n" + itemsFound + "\n" + urlsVisited + "```"
	return output
}
