package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/urfave/cli"
)

var outputDir string

func main() {
	var app = cli.NewApp()
	// set info
	app.Name = "w3kollector CLI"
	app.Usage = "scrape and scan websites"
	app.Author = "Xenophonsec"
	app.Version = "1.0.0"

	app.Commands = []cli.Command{
		{
			Name:    "scrape",
			Aliases: []string{"s"},
			Usage:   "crawl and scrape a website",
			Action: func(c *cli.Context) error {
				if c.NArg() > 0 {
					scrapeSite(c.Args().First())
				} else {
					fmt.Println("No target URL given")
				}
				return nil
			},
		},
		{
			Name:    "lookup",
			Aliases: []string{"l"},
			Usage:   "lookup domain name info, records and IP addresses",
			Action: func(c *cli.Context) error {
				if c.NArg() > 0 {
					lookupDomain(c.Args().First())
				} else {
					fmt.Println("No target domain given")
				}
				return nil
			},
		},
	}

	showBanner()

	// run cli
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func lookupDomain(domain string) {
	fmt.Println()
	lookupCNAME(domain)
	fmt.Println()
	getIPs(domain)
	fmt.Println()
	lookupNS(domain)
	fmt.Println()
	lookupMX(domain)
	fmt.Println()
	lookupDNSTXT(domain)
}

// Crawl and scrape website for email addresses, phone numbers, external links, and download pdfs and other interesting files
func scrapeSite(targetURL string) {
	c := colly.NewCollector(
		colly.IgnoreRobotsTxt(),
	)
	var emailAddresses []string
	var pdfs []string
	var phoneNumbers []string
	var pages []string
	var externalLinks []string

	if !strings.HasPrefix(targetURL, "http") {
		targetURL = "https://" + targetURL
	}

	// set and clean target domain
	targetDomain := strings.Replace(targetURL, "https://", "", -1)
	targetDomain = strings.Replace(targetDomain, "http://", "", -1)

	var err error
	outputDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	outputDir = filepath.Join(outputDir, targetDomain)
	_, err = os.Stat(outputDir)
	if os.IsNotExist(err) {
		os.Mkdir(outputDir, os.FileMode(0644))
	}

	// for every link
	c.OnHTML("a", func(e *colly.HTMLElement) {
		url := e.Attr("href")

		// handle mailto links
		if strings.Contains(url, "mailto:") {
			emailAddress := strings.Replace(url, "mailto:", "", -1)
			emailAddress = strings.Trim(emailAddress, " ")
			if !arrayContains(emailAddresses, emailAddress) {
				emailAddresses = append(emailAddresses, emailAddress)
				fmt.Println("Email Address: ", emailAddress)
				saveLineToFile("emailAddresses.txt", emailAddress+" : "+url)
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
			// handle pdfs
			if strings.HasSuffix(url, ".pdf") {
				if !arrayContains(pdfs, url) {
					pdfs = append(pdfs, url)
					fmt.Println("PDF: ", url)
					saveLineToFile("pdfs.txt", url)
				}
			}
			e.Request.Visit(url)
		} else {
			if !strings.HasPrefix(url, "#") && !arrayContains(externalLinks, url) && !strings.HasPrefix(url, "javascript:") {
				externalLinks = append(externalLinks, url)
				fmt.Println("Outbound: ", url)
				saveLineToFile("outbound.txt", url)
			}
		}
	})

	// for every request
	c.OnRequest(func(r *colly.Request) {
		// log url
		fmt.Println("Scanning: ", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		url := r.Request.URL.String()
		page := r.Request.URL.Path
		// save response meta data
		responseLine := strconv.Itoa(r.StatusCode) + " " + r.Headers.Get("Content-Type") + " " + url
		saveLineToFile("responses.txt", responseLine)
		// download pdfs
		if strings.HasSuffix(page, ".pdf") {
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
			// download files that aren't html pages
		} else if !strings.Contains(r.Headers.Get("Content-Type"), "html") {
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
		} else {
			// save page to sitemap
			if !arrayContains(pages, page) {
				pages = append(pages, page)
				saveLineToFile("sitemap.txt", page)
			}
			// scrape phone numbers
			content := string(r.Body)
			reg := regexp.MustCompile(`[0-9]{3}-[0-9]{3}-[0-9]{4}`)
			matches := reg.FindAllString(content, -1)
			for _, phonenumber := range matches {
				if !arrayContains(phoneNumbers, phonenumber) {
					phoneNumbers = append(phoneNumbers, phonenumber)
					fmt.Println("Phone Number: ", phonenumber)
					saveLineToFile("phoneNumbers.txt", phonenumber+" : "+url)
				}
			}
		}
	})

	c.Visit(targetURL)
}
func saveLineToFile(filename string, content string) {
	filePath := filepath.Join(outputDir, filename)
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to write to file", filename)
	} else {
		_, err := f.WriteString(content + "\n")
		if err != nil {
			fmt.Println("Failed to write to file", filename)
		}
	}
	if err := f.Close(); err != nil {
		log.Fatal("error", err)
	}
}

func arrayContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func lookupDNSTXT(domain string) {
	records, err := net.LookupTXT(domain)
	if err != nil {
		log.Fatal(err)
	}
	for _, record := range records {
		fmt.Println("DNS TXT Record:", record)
	}
}

func lookupCNAME(domain string) {
	cname, err := net.LookupCNAME(domain)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("CNAME Record:", cname)
}

func getIPs(domain string) {
	ips, err := net.LookupHost(domain)
	if err != nil {
		log.Fatal(err)
	}
	for _, ip := range ips {
		fmt.Println("IP address:", ip)
	}
}

func lookupNS(domain string) {
	records, err := net.LookupNS(domain)
	if err != nil {
		log.Fatal(err)
	}
	for _, record := range records {
		fmt.Println("Name Server:", record.Host)
	}
}

func lookupMX(domain string) {
	mxs, err := net.LookupMX(domain)
	if err != nil {
		log.Fatal(err)
	}
	for _, mx := range mxs {
		fmt.Println("MX Record:", mx.Host)
	}
}

func showBanner() {
	fmt.Println()
	fmt.Println("          _____ __ __      ____          __            ")
	fmt.Println(" _      _|__  // //_/___  / / /__  _____/ /_____  _____")
	fmt.Println("| | /| / //_ </ ,< / __ \\/ / / _ \\/ ___/ __/ __ \\/ ___/")
	fmt.Println("| |/ |/ /__/ / /| / /_/ / / /  __/ /__/ /_/ /_/ / /    ")
	fmt.Println("|__/|__/____/_/ |_\\____/_/_/\\___/\\___/\\__/\\____/_/     ")
	fmt.Println("=================================================")
	fmt.Println()
}
