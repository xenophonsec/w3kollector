package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/gocolly/colly"
	"github.com/urfave/cli"
)

var searches []string

func searchSite(cl *cli.Context, targetURL string, crawl bool) {
	collector := colly.NewCollector(
		colly.IgnoreRobotsTxt(),
	)
	targetURL = handleTargetURL(targetURL)
	targetDomain := urlToDomain(targetURL)

	handleOutputPath(cl.String("out"), targetDomain)
	// crawl website
	if crawl {
		collector.OnHTML("a", func(e *colly.HTMLElement) {
			url := e.Attr("href")
			if strings.Contains(url, targetDomain) ||
				// or is a relative link
				strings.HasPrefix(url, "./") ||
				( // or isn't the following
				// an external url
				!strings.HasPrefix(url, "http") &&
					// a javascript url
					!strings.HasPrefix(url, "javascript:") &&
					!strings.HasPrefix(url, "mailto:") &&
					!strings.HasPrefix(url, "tel:")) {
				// visit next page
				e.Request.Visit(url)
			}
		})
	}
	collector.OnRequest(func(r *colly.Request) {
		// log url
		fmt.Println("Searching: ", r.URL)
	})

	// search page
	collector.OnResponse(func(r *colly.Response) {

		url := r.Request.URL.String()
		content := strings.ToLower(string(r.Body))
		search := strings.ToLower(cl.String("text"))
		keywords := strings.ToLower(cl.String("keywords"))
		searchPage(url, content, search, keywords)
	})

	collector.Visit(targetURL)
}

func searchPage(url string, content string, search string, keywords string) {
	if len(keywords) > 0 {
		keys := strings.Split(keywords, ",")
		for _, k := range keys {
			occurences := strings.Count(content, k)
			if occurences >= 1 {
				if !arrayContains(searches, url) {
					searches = append(searches, url)
					output := k + " found " + strconv.Itoa(occurences) + " times at: " + url
					color.Green(output)
					saveLineToFile("searches.txt", output)
				}
			}
		}
	} else if len(search) > 0 {
		occurences := strings.Count(content, search)
		if occurences >= 1 {
			if !arrayContains(searches, url) {
				searches = append(searches, url)
				output := search + " found " + strconv.Itoa(occurences) + " times at: " + url
				color.Green(output)
				saveLineToFile("searches.txt", output)
			}
		}
	} else {
		color.Red("No search critera was specified")
		color.Red("Please specify --text or --keywords")
	}
}
