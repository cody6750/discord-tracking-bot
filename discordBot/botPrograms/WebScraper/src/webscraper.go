package webscraper

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

//Run ..

func connectToURL(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(resp.Body)
	// Convert HTML into goquery document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(doc)
}

//Run ..
func Run() {
	log.Println("Webscraper intializing.......")
	connectToURL("https://amazon.com")
}
