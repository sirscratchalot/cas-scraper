package scrapers

import (
	csv "cas-scraper/csv"
)

type Scraper interface {
	CrawlScrape(line csv.Line, ch chan Result, scraperIndex int)
}
type Result struct {
	Error       error    //Any error triggered during processing
	Line        csv.Line //Original Line if no result, line with new values added as columns otherwise
	ScrapeIndex int      //Index of the scraper used to produce this result, allowing retry with different scraper if not succesful.
}
