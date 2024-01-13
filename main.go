// This program generates a words.txt file for the game Wordle
// It does so by scraping the site: 0
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

const (
	url         = "https://www.wordunscrambler.net/word-list/wordle-word-list"
	mainDomain  = "www.wordunscrambler.net"
	myUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2.1 Safari/605.1.15"
)

var (
	words strings.Builder
)

func main() {
	c := colly.NewCollector(
		colly.UserAgent(myUserAgent),
		colly.AllowedDomains(mainDomain),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Fatalln("Something went wrong: ", err)
	})

	c.OnHTML("div[class=content]", func(e *colly.HTMLElement) {
		e.ForEach("div", func(i int, e *colly.HTMLElement) {
			elemText := e.ChildText("h3")
			if strings.HasPrefix(elemText, "Wordle Words List Starting With") {
				e.ForEach("li a", func(i int, h *colly.HTMLElement) {
					word := strings.TrimSpace(h.Text)
					words.WriteString(word)
					words.WriteByte('\n')
				})
			}
		})
	})

	c.OnScraped(func(r *colly.Response) {
		fname := "words.txt"
		err := os.WriteFile(fname, []byte(words.String()), os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("File '%s' written to disk\n", fname)
	})

	c.Visit(url)
}
