// This program generates a words.txt file for the game Wordle
// It does so by scraping the site: 0
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
	"github.com/spf13/viper"
)

const (
	url        = "https://www.wordunscrambler.net/word-list/wordle-word-list"
	mainDomain = "www.wordunscrambler.net"
)

var (
	words          strings.Builder
	numWords       int
	collyUserAgent string
)

func init() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: couldn't read config file (%s)\n", err)
		log.Println("Using default UserAgent")
		return
	}
	collyUserAgent = viper.GetString("COLLY_USER_AGENT")
	log.Printf("Using UserAgent from %s\n", viper.ConfigFileUsed())
}

func main() {
	c := colly.NewCollector(
		colly.UserAgent(collyUserAgent),
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
					numWords++
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
		fmt.Printf("File '%s' written to disk (%d entries)\n", fname, numWords)
	})

	c.Visit(url)
}
