# MGoBlog-Web-Scraper
A tool for scraping article content from MGoBlog

## Usage
Scrape content from MGoBlog
```
$ ./MGoBlog-Web-Scraper -scrape=true
```
Specify the earliest date to scrape articles until (articles are retrieved in reverse chronological order. Default is today - 30 days)
```
$ ./MGoBlog-Web-Scraper -scrape=true -date=2015-09-17
```
Process all articles that have been cached locally and print out the results
```
$ ./MGoBlog-Web-Scraper -process=true
```
Specify the number of results to print after processing articles (Default is 20)
```
$ ./MGoBlog-Web-Scraper -process=true -results=30
```

## TODO
The plan is to eventually do some sort of language analysis with the data that I've collected (e.g. generate flesch-kincaid or coleman-liau scores)
