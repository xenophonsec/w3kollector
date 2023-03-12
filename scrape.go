package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/gocolly/colly"
	"github.com/urfave/cli"
)

// Crawl and scrape website for email addresses, phone numbers, external links, and download pdfs and other interesting files
func scrape(cl *cli.Context, targetURL string, crawl bool) {
	collector := colly.NewCollector(
		colly.IgnoreRobotsTxt(),
	)
	var emailAddresses []string
	var pdfs []string
	var phoneNumbers []string
	var pages []string
	var externalLinks []string
	var scripts []string
	var stylesheets []string
	var links []string
	var metaData []string
	var serverEngines []string
	var addresses []string
	var images []string

	targetURL = handleTargetURL(targetURL)
	targetDomain := urlToDomain(targetURL)

	handleOutputPath(cl.String("out"), targetDomain)
	// for every link
	collector.OnHTML("a", func(e *colly.HTMLElement) {
		url := e.Attr("href")

		// handle mailto links
		if strings.Contains(url, "mailto:") {
			if cl.Bool("emails") || cl.Bool("all") {
				emailAddress := strings.Replace(url, "mailto:", "", -1)
				emailAddress = strings.Trim(emailAddress, " ")
				if !arrayContains(emailAddresses, emailAddress) {
					emailAddresses = append(emailAddresses, emailAddress)
					color.Cyan("Email Address: " + emailAddress)
					saveLineToFile("emailAddresses.txt", emailAddress)
				}
			}
			// if the url contains the target domain
		} else if strings.Contains(url, targetDomain) ||
			// or is a relative link
			strings.HasPrefix(url, "./") ||
			( // or isn't the following
			// an external url
			!strings.HasPrefix(url, "http") &&
				// a javascript url
				!strings.HasPrefix(url, "javascript:")) {
			if strings.HasPrefix(url, "tel:") {
				if cl.Bool("phones") || cl.Bool("all") {
					phonenumber := strings.TrimPrefix(url, "tel:")
					if !arrayContains(phoneNumbers, phonenumber) {
						phoneNumbers = append(phoneNumbers, phonenumber)
						color.Magenta("Phone Number: " + phonenumber)
						saveLineToFile("phoneNumbers.txt", phonenumber)
					}
				}
			} else {
				// handle pdfs
				if strings.HasSuffix(url, ".pdf") {
					if cl.Bool("logpdfs") && !arrayContains(pdfs, url) {
						pdfs = append(pdfs, url)
						color.Blue("PDF: " + url)
						saveLineToFile("pdfs.txt", url)
					}
				}
				// only visit url if crawl is set to true or it is a pdf
				if crawl || strings.HasSuffix(url, ".pdf") {
					e.Request.Visit(url)
				}
			}
		} else {

			if !strings.HasPrefix(url, "#") && !arrayContains(externalLinks, url) && !strings.HasPrefix(url, "javascript:") {
				externalLinks = append(externalLinks, url)
				color.Red("Outbound: " + url)
				saveLineToFile("outbound.txt", url)
			}
		}
	})

	collector.OnHTML("script", func(e *colly.HTMLElement) {
		if cl.Bool("scripts") || cl.Bool("all") {
			url := e.Attr("src")
			if len(url) > 1 && !arrayContains(scripts, url) {
				scripts = append(scripts, url)
				fmt.Println("Script: ", url)
				saveLineToFile("scripts.txt", url)
			}
		}
	})

	collector.OnHTML("img", func(e *colly.HTMLElement) {
		if cl.Bool("images") || cl.Bool("all") {
			url := e.Attr("src")
			if len(url) > 1 && !arrayContains(images, url) {
				images = append(images, url)
				fmt.Println("Image: ", url)
				saveLineToFile("images.txt", url)
				e.Request.Visit(url)
			}
		}
	})

	collector.OnHTML("link", func(e *colly.HTMLElement) {
		if cl.Bool("resourcelinks") || cl.Bool("all") {
			url := e.Attr("href")
			rel := e.Attr("rel")
			if rel == "stylesheet" {
				if !arrayContains(stylesheets, url) {
					stylesheets = append(stylesheets, url)
					fmt.Println("Stylesheet: ", url)
					saveLineToFile("stylesheets.txt", url)
				}
			} else if !arrayContains(links, url) {
				links = append(links, url)
				fmt.Println("Resource Link: ", rel+": "+url)
				saveLineToFile("resourceLinks.txt", rel+": "+url)
			}
		}
	})

	collector.OnHTML("meta", func(e *colly.HTMLElement) {
		if cl.Bool("meta") || cl.Bool("all") {
			name := e.Attr("name")
			property := e.Attr("property")
			content := e.Attr("content")
			value := e.Attr("value")
			meta := name + property
			if len(content) <= 0 {
				meta += " : " + value
			} else {
				meta += " : " + content
			}
			// make sure that something we care about is set
			if len(name+property+content+value) != 0 {
				// make sure we haven't already recorded this
				if !arrayContains(metaData, meta) {
					metaData = append(metaData, meta)
					color.Yellow("Meta Tag: " + meta)
					saveLineToFile("metaTags.txt", meta)
				}
			}
		}
	})

	// for every request
	collector.OnRequest(func(r *colly.Request) {
		// log url
		fmt.Println("Scanning: ", r.URL)
	})

	collector.OnResponse(func(r *colly.Response) {
		url := r.Request.URL.String()
		page := r.Request.URL.Path
		// save response meta data
		responseLine := r.Headers.Get("server") + "\t" +
			r.Headers.Get("via") + "\t" +
			strconv.Itoa(r.StatusCode) + "\t" +
			r.Headers.Get("Content-Type") + "\t" +
			url

		saveLineToFile("responses.txt", responseLine)
		// Get server engine names
		if cl.Bool("serverEngine") || cl.Bool("all") {
			engine := r.Headers.Get("server")
			if len(engine) > 0 && !arrayContains(serverEngines, engine) {
				serverEngines = append(serverEngines, engine)
				color.Yellow("Server Engine: " + engine)
				saveLineToFile("ServerEngines.txt", engine)
			}
		}
		// download pdfs
		if strings.HasSuffix(page, ".pdf") {
			if !dontsave {
				if cl.Bool("downloadpdfs") || cl.Bool("all") {
					pdfsDir := filepath.Join(outputDir, "pdfs")
					_, err := os.Stat(pdfsDir)
					if os.IsNotExist(err) {
						os.Mkdir(pdfsDir, os.FileMode(0644))
					}
					pdfFilePath := filepath.Join(pdfsDir, r.FileName())
					_, err = os.Stat(pdfFilePath)
					if os.IsNotExist(err) {
						err = r.Save(pdfFilePath)
						if err != nil {
							panic(err)
						}
					}
				}
			}
			// download files that aren't html pages
		} else if !strings.Contains(r.Headers.Get("Content-Type"), "html") {
			if !dontsave {
				if cl.Bool("downloadpdfs") || cl.Bool("all") {
					filesDir := filepath.Join(outputDir, "files")
					_, err := os.Stat(filesDir)
					if os.IsNotExist(err) {
						os.Mkdir(filesDir, os.FileMode(0644))
					}
					filePath := filepath.Join(filesDir, r.FileName())
					_, err = os.Stat(filePath)
					if os.IsNotExist(err) {
						err = r.Save(filePath)
						if err != nil {
							panic(err)
						}
					}
				}
			}
			// if it is an HTML page
		} else {
			content := string(r.Body)
			// save page to sitemap
			if !arrayContains(pages, page) {
				pages = append(pages, page)
				saveLineToFile("sitemap.txt", page)
			}
			// handle search
			search := cl.String("search")
			if search != "" {
				occurences := strings.Count(content, search)
				if occurences >= 1 {
					output := search + " found " + strconv.Itoa(occurences) + " times at: " + url
					color.Green(output)
					saveLineToFile("search.txt", output)
				}
			}
			// scrape addresses
			if cl.Bool("addresses") || cl.Bool("all") {
				tags := regexp.MustCompile(`<[^>]*>`)
				contentStripped := tags.ReplaceAllString(content, " ")
				// I want multiple more specific matching patterns because addresses can be so vague you can easily grab something that isn't one
				// match P.O. Box addresses
				// P.O. Box 3000 Cityname, CA 12345
				findAddresses(contentStripped, `(P\.O\.\sBox\s)\d+\s[\w+\s]+[,]?\s\w+\s\w+`, &addresses)
				// standard US home address
				// 889 easy street rd, happyville Virginia 22642
				findAddresses(contentStripped, `\d+\s+[\w+\s]+\s[crdstav]{2}[,]?\s\w+\s\w+\s\d{3,5}`, &addresses)
				// building floors
				// 200 Broadway 4th Floor New York, NY 10003
				findAddresses(contentStripped, `\d+\s+[\w+\s]+floor[,]?\s\w+\s\w+[,]?\s\w+\s\d{3,5}`, &addresses)
				// suites and apt
				//260 Deanna Views, Suite 480, Robynberg, Montana, 56479
				//260 Deanna Views, Apt. 480, Robynberg, Montana, 56479
				findAddresses(contentStripped, `\d+\s+[\w+\s]+[,]?\s+(Apt\.|Suite)\s\d+[,]?\s+\w+[,]?\s+\w+[,]?\s+\d+`, &addresses)
			}
			// scrape phone numbers
			if cl.Bool("phones") || cl.Bool("all") {
				reg := regexp.MustCompile(`[0-9]{3}-[0-9]{3}-[0-9]{4}`)
				matches := reg.FindAllString(content, -1)
				for _, phonenumber := range matches {
					if !arrayContains(phoneNumbers, phonenumber) {
						phoneNumbers = append(phoneNumbers, phonenumber)
						color.Magenta("Phone Number: " + phonenumber)
						saveLineToFile("phoneNumbers.txt", phonenumber)
					}
				}
				reg = regexp.MustCompile(`\([0-9]{3}\)\s{0,1}[0-9]{3}-[0-9]{4}`)
				matches = reg.FindAllString(content, -1)
				for _, phonenumber := range matches {
					if !arrayContains(phoneNumbers, phonenumber) {
						phoneNumbers = append(phoneNumbers, phonenumber)
						color.Magenta("Phone Number: " + phonenumber)
						saveLineToFile("phoneNumbers.txt", phonenumber)
					}
				}
			}
		}
	})

	collector.Visit(targetURL)
}

func findAddresses(text string, regex string, addresses *[]string) {
	reg := regexp.MustCompile(regex)
	matches := reg.FindAllString(text, -1)
	for _, address := range matches {
		// clean address whitespace
		spaces := regexp.MustCompile(`\s+`)
		address = spaces.ReplaceAllString(address, " ")
		if !arrayContains(*addresses, address) {
			*addresses = append(*addresses, address)
			color.Cyan("Possible Mailing Address: " + address)
			saveLineToFile("addresses.txt", address)
		}
	}
}
