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

var (
	currentMetrics Metrics
	totalMetrics   Metrics
)

func EnableCurrentMetrics(m Metrics) Metrics {
	currentMetrics = m
	return currentMetrics
}

func EnableTotalMetrics(m Metrics) Metrics {
	totalMetrics = m
	return totalMetrics
}

func generateMetricsOutput(m *Metrics) string {
	Header := "Scraping Metrics :" + m.URL
	duplicateUrlsFound := "Duplicate Urls Found: " + strconv.Itoa(m.DuplicatedUrlsFound)
	urlsFound := "Urls Found: " + strconv.Itoa(m.UrlsFound)
	itemsFound := "Items Found: " + strconv.Itoa(m.ItemsFound)
	urlsVisited := "Urls visited Found: " + strconv.Itoa(m.ItemsFound)
	output := "\n```" + Header + "\n" + duplicateUrlsFound + "\n" + urlsFound + "\n" + itemsFound + "\n" + urlsVisited + "```"
	return output
}
