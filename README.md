# MGoBlog-Web-Scraper
A tool for scraping article content from MGoBlog

## Usage
Build from source and then run the executable with one of the following command line flags
* -scrape=true
  * -date=<YYYY-MM-DD>
* -process=true
  * -number=numberToPrint

If scrape is true then all articles after <date> will be scraped from the website and cached locally on disk. If the date argument is not supplied the default date is used which at this moment is 2015-07-01. If process is set to true then all locally cached files will be parsed and word count of the numberToPrint most common words aggregated from all rticles will be printed out.

## TODO
The plan is to eventually do some sort of language analysis with the data that I've collected (e.g. generate flesch-kincaid or coleman-liau scores)
