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
	//currentMetrics represents the metrics recieved from the current REST API call.
	currentMetrics Metrics

	//totalMetrics represents the metrics recieved from the all of the REST API calls.
	totalMetrics Metrics
)

//SetCurrentMetrics used to pass in the current metrics from the trackingbot package to avoid circular dependency
func SetCurrentMetrics(m Metrics) Metrics {
	currentMetrics = m
	return currentMetrics
}

//SetTotalMetrics used to pass in the total metrics from the trackingbot package to avoid circular dependency
func SetTotalMetrics(m Metrics) Metrics {
	totalMetrics = m
	return totalMetrics
}

//AppendMetrics adds two metric sturcts.
func AppendMetrics(m, a *Metrics) {
	m.DuplicatedUrlsFound += a.DuplicatedUrlsFound
	m.ItemsFound += a.ItemsFound
	m.UrlsFound += a.UrlsFound
	m.UrlsVisited += a.UrlsVisited
}

//GenerateMetricsOutput generates formatted outputs given a metrics struct.
func GenerateMetricsOutput(m Metrics) string {
	Header := "Scraping Metrics :" + m.URL
	duplicateUrlsFound := "Duplicate Urls Found: " + strconv.Itoa(m.DuplicatedUrlsFound)
	urlsFound := "Urls Found: " + strconv.Itoa(m.UrlsFound)
	itemsFound := "Items Found: " + strconv.Itoa(m.ItemsFound)
	urlsVisited := "Urls visited Found: " + strconv.Itoa(m.ItemsFound)
	output := "\n```" + Header + "\n" + duplicateUrlsFound + "\n" + urlsFound + "\n" + itemsFound + "\n" + urlsVisited + "```"
	return output
}
