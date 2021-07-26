package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

type profile struct {
	Name      string
	Image     string
	Title     string
	Bio       string
	BirthDate string
	TopMovies []movie
}

type movie struct {
	Title         string
	CharacterName string
	Year          string
}

func main() {
	month := flag.Int("month", 7, "Month to fetch")
	day := flag.Int("day", 30, "Day to fetch")
	flag.Parse()
	crawl(*month, *day)
	fmt.Println("Hello World")
}

func crawl(month int, day int) {
	collector := colly.NewCollector(
		colly.AllowedDomains("www.imdb.com", "imdb.com"),
	)

	infoCollector := collector.Clone()

	collector.OnHTML(".mode-detail", func(h *colly.HTMLElement) {
		profileURL := h.ChildAttr("div.lister-item-image > a", "href")
		profileURL = h.Request.AbsoluteURL(profileURL)

		infoCollector.Visit(profileURL)
	})

	collector.OnHTML("a.lister-page-next", func(h *colly.HTMLElement) {
		nextPage := h.Request.AbsoluteURL(h.Attr("href"))
		collector.Visit(nextPage)
	})

	infoCollector.OnHTML("#content-2-wide", func(h *colly.HTMLElement) {
		tempProfile := profile{}
		tempProfile.Name = h.ChildText("h1.header > span.itemprop")
		fmt.Println("Getting data for: ", tempProfile.Name)
		tempProfile.Image = h.ChildAttr("#name-poster", "src")
		tempProfile.Title = h.ChildText("#name-job-categories > a > span.itemprop")
		tempProfile.Bio = strings.TrimSpace(h.ChildText("#name-bio-text > div.name-trivia-bio-text > div.inline"))
		tempProfile.BirthDate = h.ChildAttr("#name-born-info time", "datetime")

		h.ForEach("div.knownfor-title", func(i int, m *colly.HTMLElement) {
			tempMovie := movie{}
			tempMovie.Title = m.ChildText("div.knownfor-title-role > a")
			tempMovie.CharacterName = m.ChildText("div.knownfor-title-role > span")
			tempMovie.Year = m.ChildText("div.knownfor-year > span")

			tempProfile.TopMovies = append(tempProfile.TopMovies, tempMovie)
		})

		jsonRes, err := json.MarshalIndent(tempProfile, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(jsonRes))
	})

	// collector.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting: ", r.URL.String())
	// })

	// infoCollector.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting: ", r.URL.String())
	// })

	baseurl := fmt.Sprintf("https://www.imdb.com/search/name/?birth_monthday=%d-%d", month, day)
	collector.Visit(baseurl)
}
