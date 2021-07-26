package main

// Importing the required packages
import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// Creating user-defined type "profile"
type profile struct {
	Name      string
	Image     string
	Title     string
	Bio       string
	BirthDate string
	TopMovies []movie
}

// Creating user-defined type "movie"
type movie struct {
	Title         string
	CharacterName string
	Year          string
}

// main function
func main() {

	// Getting user input for Month
	fmt.Println("Enter the Month code (Eg. 1 for Jan) :")
	var month int
	fmt.Scanln(&month)

	// Getting user input for Day
	fmt.Println("Enter the Day :")
	var day int
	fmt.Scanln(&day)

	// Getting user input for Number of Profiles
	fmt.Println("Number of profiles to fetch (Default is 5) :")
	var numberOfProfile int
	fmt.Scanln(&numberOfProfile)

	fmt.Println("Started to crawl data ....")

	// Calling the crawl function to start crawling
	crawl(month, day, numberOfProfile)
}

func crawl(month int, day int, numberOfProfile int) {

	// Creating a slice of type profile
	finalData := []profile{}

	// Initializing our counter to zero
	count := 0

	// Setting the default value of numberOfProfile as 5
	if numberOfProfile == 0 {
		numberOfProfile = 5
	}

	// Creating a new Collector instance with default configuration
	collector := colly.NewCollector(
		colly.AllowedDomains("www.imdb.com", "imdb.com"),
	)

	// Creating another collector infoCollector
	infoCollector := collector.Clone()

	// Getting the profileURL and make a request to that url
	collector.OnHTML(".mode-detail", func(h *colly.HTMLElement) {
		profileURL := h.ChildAttr("div.lister-item-image > a", "href")
		profileURL = h.Request.AbsoluteURL(profileURL)

		if count < numberOfProfile {
			infoCollector.Visit(profileURL)
		}

	})

	// Getting profiledata and append that to finalData
	infoCollector.OnHTML("#content-2-wide", func(h *colly.HTMLElement) {
		tempProfile := profile{}
		tempProfile.Name = h.ChildText("h1.header > span.itemprop")
		count++
		fmt.Println(strconv.Itoa(count), "> Getting data for: ", tempProfile.Name)
		tempProfile.Image = h.ChildAttr("#name-poster", "src")
		tempProfile.Title = h.ChildText("#name-job-categories > a > span.itemprop")
		tempProfile.Bio = strings.TrimSpace(h.ChildText("#name-bio-text > div.name-trivia-bio-text > div.inline"))
		tempProfile.Bio = strings.TrimSpace(strings.ReplaceAll(tempProfile.Bio, "See full bio »", ""))
		tempProfile.BirthDate = h.ChildAttr("#name-born-info time", "datetime")

		h.ForEach("div.knownfor-title", func(i int, m *colly.HTMLElement) {
			tempMovie := movie{}
			tempMovie.Title = m.ChildText("div.knownfor-title-role > a")
			tempMovie.CharacterName = m.ChildText("div.knownfor-title-role > span")
			tempMovie.Year = m.ChildText("div.knownfor-year > span")

			tempProfile.TopMovies = append(tempProfile.TopMovies, tempMovie)
		})

		// Appending each profile to finalData
		finalData = append(finalData, tempProfile)

	})

	// collector.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting: ", r.URL.String())
	// })

	// infoCollector.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting: ", r.URL.String())
	// })

	// Making request to the baseurl
	baseurl := fmt.Sprintf("https://www.imdb.com/search/name/?birth_monthday=%d-%d", month, day)
	collector.Visit(baseurl)

	jsonRes, err := json.MarshalIndent(finalData, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(string(jsonRes))

	// Setting the filename, Ex: "1-1data.json" if month = 1 and day = 1
	fileName := fmt.Sprintf("%d-%ddata.json", month, day)

	// Calling the writeFile function to write the data into the json file
	writeFile(fileName, string(jsonRes))

	fmt.Println("Done !")

}

// Function for writing the json file with the data
func writeFile(fileName string, data string) {
	f, openerr := os.Create(fileName)

	if openerr != nil {
		log.Fatal(openerr)
	}

	defer f.Close()

	_, err2 := f.WriteString(data)

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println("Writing file:", fileName)
}
