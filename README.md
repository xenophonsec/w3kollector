![w3Kollector banner](./_images/w3kollector_banner.png "Don't do what you don't have to")
# w3Kollector
A greedy website scanner, scraper, and crawler.

Get everything from email addresses and phone numbers to PDFs and DNS records.

## Install

```
go install -v github.com/xenophonsec/w3kollector@latest
```

## Use

```
          _____ __ __      ____          __
 _      _|__  // //_/___  / / /__  _____/ /_____  _____
| | /| / //_ </ ,< / __ \/ / / _ \/ ___/ __/ __ \/ ___/
| |/ |/ /__/ / /| / /_/ / / /  __/ /__/ /_/ /_/ / /
|__/|__/____/_/ |_\____/_/_/\___/\___/\__/\____/_/
=================================================

NAME:
   w3kollector - scrape and scan websites

USAGE:
   w3kollector [global options] command [command options] [arguments...]

VERSION:
   1.0.0

AUTHOR:
   Xenophonsec

COMMANDS:
   lookup, l  lookup domain name info, records and IP addresses
   scrape, s  scrape html
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## Features

### Lookup

The lookup command fetches

- CNAME
- IP addresses
- Name servers
- MX records
- DNS TXT Records

```
w3kollector lookup microsoft.com
```
![w3Kollector lookup](./_images/w3kollector_lookup.png)

### Scrape

The scraper crawls websites and gathers contact info, meta data and various files. It refrains from following links outside of the target domain (but it will scan subdomains).

```
w3kollector scrape --help
```

```
          _____ __ __      ____          __
 _      _|__  // //_/___  / / /__  _____/ /_____  _____
| | /| / //_ </ ,< / __ \/ / / _ \/ ___/ __/ __ \/ ___/
| |/ |/ /__/ / /| / /_/ / / /  __/ /__/ /_/ /_/ / /
|__/|__/____/_/ |_\____/_/_/\___/\___/\__/\____/_/
=================================================

NAME:
   w3kollector scrape - scrape html

USAGE:
   w3kollector scrape [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
   crawl  crawl the website
   page   scrape a specific page

GLOBAL OPTIONS:
   --help, -h  show help
```

You can scrape a whole website with `scrape crawl` or a single page with the `scrape page` command.

```
w3kollector scrape crawl --help
```
```
          _____ __ __      ____          __
 _      _|__  // //_/___  / / /__  _____/ /_____  _____
| | /| / //_ </ ,< / __ \/ / / _ \/ ___/ __/ __ \/ ___/
| |/ |/ /__/ / /| / /_/ / / /  __/ /__/ /_/ /_/ / /
|__/|__/____/_/ |_\____/_/_/\___/\___/\__/\____/_/
=================================================

NAME:
   w3kollector scrape crawl - crawl the website

USAGE:
   w3kollector scrape crawl [command options] [arguments...]

OPTIONS:
   --all, -a                     Scrape everything
   --emails, -e                  Scrape email addresses
   --phones, -p                  Scrape phone numbers
   --scripts, -s                 Scrape scripts
   --stylesheets, --ss           Scrape stylesheets
   --resourcelinks, --rl         Scrape resource links
   --downloadpdfs, --dp          Download pdfs
   --files                       Download all downloadable files
   --logpdfs, --lp               Log pdf urls
   --meta, -m                    Scrape meta tags
   --search value, --find value  Search html pages. This can be plain text or html you are looking for.

```


```
w3kollector scrape crawl --all microsoft.com
```
![w3Kollector scrape](./_images/w3kollector_scrape.png)

It collects
- Email Addresses
- Phone Numbers
- PDFs
- Downloadable files
- Outbound Links
- Scripts
- Stylesheets
- Resource Links (dns-prefetch, preconnect, canonical, alternate etc...)
- Meta data from meta tags

It also
- builds a sitemap of pages visited
- logs http responses
- supports raw html text search

w3Kollector will write files to whatever directory you are currently in.

Output structure:
- website.com
  - files
    - setup.exe
  - pdfs
    - gettingstarted.pdf
  - outbound.txt
  - sitemap.txt
  - responses.txt
  - emailAddresses.txt
  - phoneNumbers.txt
  - scripts.txt
  - stylesheets.txt
  - resourceLinks.txt
  - metaTags.txt
  - search.txt

> **Note:** files and folders are only created if that data is found

Resource links are stored with both the href and the rel attributes:
```
preconnect: https://fonts.gstatic.com
icon: /favicon.png
canonical: https://example.com/blog/
alternate: https://example.com/rss/
```

> **Note:** Only unique values are stored. If it finds the same phone number somewhere else it will be skipped.

If you provide a domain name and not a full url it will auto prefix https://
If you want to specify http you can do so likewise:
```
w3kollector scrape crawl --all http://website.com
```
