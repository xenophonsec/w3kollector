package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

var outputDir string
var dontsave bool

func main() {
	var app = cli.NewApp()
	// set info
	app.Name = "w3kollector"
	app.Usage = "scrape and scan websites"
	app.Author = "Xenophonsec"
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "dontsave,ds",
			Usage:       "Do not save results to files",
			Destination: &dontsave,
		},
	}

	app.Commands = []cli.Command{
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
		{
			Name:    "scrape",
			Aliases: []string{"s"},
			Usage:   "scrape html",
			Subcommands: []cli.Command{
				{
					Name:  "site",
					Usage: "crawl the website and scrape it",
					Flags: scrapeFlags,
					Action: func(c *cli.Context) error {
						if c.NArg() > 0 {
							scrape(c, c.Args().First(), true)
						} else {
							fmt.Println("No target URL given")
						}
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
				{
					Name:  "page",
					Usage: "scrape a specific page",
					Flags: scrapeFlags,
					Action: func(c *cli.Context) error {
						if c.NArg() > 0 {
							scrape(c, c.Args().First(), false)
						} else {
							fmt.Println("No target URL given")
						}
						return nil
					},
				},
			},
		},
		{
			Name:    "search",
			Aliases: []string{"q"},
			Usage:   "Search website for specific text or html",
			Subcommands: []cli.Command{
				{
					Name:  "site",
					Usage: "crawl and search the whole website",
					Flags: searchFlags,
					Action: func(c *cli.Context) error {
						if c.NArg() > 0 {
							searchSite(c, c.Args().First(), true)
						} else {
							fmt.Println("No target URL given")
						}
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
				{
					Name:  "page",
					Usage: "search a specific page",
					Flags: searchFlags,
					Action: func(c *cli.Context) error {
						if c.NArg() > 0 {
							searchSite(c, c.Args().First(), false)
						} else {
							fmt.Println("No target URL given")
						}
						return nil
					},
				},
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
