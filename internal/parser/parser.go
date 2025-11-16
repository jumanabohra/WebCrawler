package parser

import (
	"bytes"
	"log"

	"github.com/PuerkitoBio/goquery"
)

type Parser struct {
}

func New() *Parser {
	return &Parser{}
}

// Parse the body and return list of links found in the body
func (p *Parser) Parse(body []byte) ([]string, error) {

	links := []string{}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatalf("Failed to parse HTML body: %v", err)
		return nil, err
	}

	doc.Find("a").Each(func(index int, item *goquery.Selection) {
		link, _ := item.Attr("href")
		links = append(links, link)
	})
	return links, nil
}
