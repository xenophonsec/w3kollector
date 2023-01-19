package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

var outputDir string
var recursive string

func main() {
	var app = cli.NewApp()
	// set info
	app.Name = "w3kollector CLI"
	app.Usage = "scrape and scan websites"
	app.Author = "Xenophonsec"
	app.Version = "1.0.0"

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
					Name:  "crawl",
					Usage: "crawl the website",
					Flags: scrapeFlags,
					Action: func(c *cli.Context) error {
						if c.NArg() > 0 {
							scrape(c.Args().First(), true)
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
							scrape(c.Args().First(), false)
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
