package main

import (
	"github.com/urfave/cli"
)

var scrapeFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "all,a",
		Usage: "Scrape everything",
	},
	cli.BoolFlag{
		Name:  "emails,e",
		Usage: "Scrape email addresses",
	},
	cli.BoolFlag{
		Name:  "phones,p",
		Usage: "Scrape phone numbers",
	},
	cli.BoolFlag{
		Name:  "scripts,s",
		Usage: "Scrape scripts",
	},
	cli.BoolFlag{
		Name:  "stylesheets,ss",
		Usage: "Scrape stylesheets",
	},
	cli.BoolFlag{
		Name:  "resourcelinks,rl",
		Usage: "Scrape resource links",
	},
	cli.BoolFlag{
		Name:  "downloadpdfs,dp",
		Usage: "Download pdfs",
	},
	cli.BoolFlag{
		Name:  "files",
		Usage: "Download all downloadable files",
	},
	cli.BoolFlag{
		Name:  "logpdfs,lp",
		Usage: "Log pdf urls",
	},
	cli.BoolFlag{
		Name:  "meta,m",
		Usage: "Scrape meta tags",
	},
	cli.StringFlag{
		Name:  "search,find,lookfor",
		Usage: "Search html pages. This can be plain text or html you are looking for",
	},
	cli.StringFlag{
		Name:  "out,o",
		Value: "",
		Usage: "What directory to place the output in. Default is current working directory",
	},
	cli.BoolFlag{
		Name:  "serverEngine,se",
		Usage: "Get server engine names from HTTP headers",
	},
}

var searchFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "text,t",
		Value: "",
		Usage: "text to search website for",
	},
	cli.StringFlag{
		Name:  "keywords,k",
		Value: "",
		Usage: "keywords seperated by commas",
	},
}
