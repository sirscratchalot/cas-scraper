//go:binary-only-package
package main

import (
	"bufio"
	csvutil "cas-scraper/csv"
	"cas-scraper/scrapers"
	"cas-scraper/scrapers/molbase"
	"cas-scraper/scrapers/webbook"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

var availableScrapers = map[string]scrapers.Scraper{
	"webbook.nist.gov": webbook.ScrapeWebbook{},
	"molbase.com":      molbase.ScrapeMolbase{},
}
var orderedScrapers = []scrapers.Scraper{
	webbook.ScrapeWebbook{},
	molbase.ScrapeMolbase{},
}

func main() {

	var input string
	var output string
	var source string
	var quick string

	setupFlags(&input, &output, &source, &quick)
	flag.Parse()

	fmt.Printf(" HELLOOO %s", os.Args)
	lines := [][]string{}
	if quick == "" && input == "" {
		fmt.Println("Please provide either -input csv file or -quick single CAS number")
		os.Exit(1)
	}
	if output == "" {
		output = "-"
	}
	if quick == "" {
		lines = readFile(input)
	} else {
		lines = [][]string{{quick}}
		fmt.Printf("Quick lines %s\n\r", lines)
	}

	headerRow := csvutil.CheckHeaderRow(lines)
	casNumberColumn := csvutil.CheckCasColumn(lines, headerRow)

	scraper := availableScrapers[source]
	//TODO: Refactor to call once for each CAS and allow cycling providers.
	outputLines := asyncCrawl(lines, headerRow, casNumberColumn, scraper)
	writeFile(outputLines, output)

}
func setupFlags(input *string, output *string, source *string, quick *string) {
	available := ""
	for k, _ := range availableScrapers {
		available = available + k + "\n"
	}
	flag.StringVar(source, "source", "", "Source site for retrieving data. If not set the scraping will use scrapers in order until procurring a result.\r\n"+available)
	flag.StringVar(input, "input", "", "Source CSV file for input. Expected to have single column or column label of 'Cas number'\nPossible values:")
	flag.StringVar(output, "output", "", "Output CSV file containing results of parsing. Will contain original info from CSV")
	flag.StringVar(quick, "quick", "", "Provide single CAS id for testing or quick reference")
}

/**
* Creates a Go routine per cas-number and awaits scraping.
* If a single scraper is provided only this will be used. Otherwise errors will trigger re-scraping with the next scraper.
 */
func asyncCrawl(lines [][]string, headerRow bool, casColumn int, scraper scrapers.Scraper) [][]string {
	targetResponses := len(lines)

	if headerRow {
		targetResponses = targetResponses - 1
	}

	channel := make(chan scrapers.Result)
	for i, line := range lines {
		//Ignore header row if found
		fmt.Printf("%d %s", i, line)
		if !headerRow || i > 0 {
			line := csvutil.Line{Columns: line, RowNumber: i, CasColumn: casColumn}
			if scraper == nil {
				go orderedScrapers[0].CrawlScrape(line, channel, 0)
			} else {
				go scraper.CrawlScrape(line, channel, 0)
			}
		}
	}
	//Waits for completion of all Go routines, resubmits go routine for errors.
	lineBuffer := make([][]string, len(lines))
	for i := 0; i < targetResponses; i++ {
		if !headerRow || i > 0 {
			result := <-channel
			lineBuffer[result.Line.RowNumber] = result.Line.Columns
			if result.Error == nil || scraper != nil || result.ScrapeIndex+1 == len(orderedScrapers) {
				fmt.Printf("Result for %s, error: %s", result.Line.Columns[casColumn], result.Error)
			} else {
				fmt.Printf("Error! Retrying with next scraper for %s", result.Line.Columns[casColumn])
				targetResponses++ //Wait for one more result.
				go orderedScrapers[0].CrawlScrape(result.Line, channel, result.ScrapeIndex+1)
			}

		}
	}
	return lineBuffer

}

func readFile(filePath string) [][]string {
	fmt.Printf("Reading file %s\n\r", filePath)
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Could not read file %s\n\r", filePath)
		os.Exit(1)
	}
	defer f.Close()
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		fmt.Printf("Could not parse CSV, is %s a CSV file?\n\r", filePath)
		os.Exit(1)
	}
	fmt.Printf("Lines %s \n\r", lines)
	return lines
}

func writeFile(lines [][]string, targetFile string) {
	/*	_, err := os.Stat(targetFile)
		if err != nil && os.IsNotExist(err) {
			os.Create(targetFile)
		}*/
	w := csv.NewWriter(os.Stdout)
	if targetFile != "-" {
		f, err := os.OpenFile(targetFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Printf("Could not write to file %s, due to: %s.\nPrinting to Stdout", targetFile, err.Error())
		} else {
			defer f.Close()
			w = csv.NewWriter(bufio.NewWriter(f))
		}
	}
	w.WriteAll(lines)
}
