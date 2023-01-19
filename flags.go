package main

import "github.com/urfave/cli"

var scrapeFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "all,a",
		Usage: "Scrape everything",
	},
	cli.BoolFlag{
		Name:  "email,e",
		Usage: "Scrape email addresses",
	},
	cli.BoolFlag{
		Name:  "phone,p",
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
		Name:  "logpdfs,lp",
		Usage: "Log pdf urls",
	},
	cli.BoolFlag{
		Name:  "meta,m",
		Usage: "Scrape meta tags",
	},
}
