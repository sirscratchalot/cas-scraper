package scrapers

type Scraper interface {
	RunScrape(inputLines [][]string, headerRow bool, casColumn int) ([][]string, error)
}
