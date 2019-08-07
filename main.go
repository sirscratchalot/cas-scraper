//go:binary-only-package
package main

import (
	"bufio"
	"cas-scraper/scrapers"
	"cas-scraper/scrapers/webbook"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"regexp"
)

var availableScrapers = map[string]scrapers.Scraper{"webbook.nist.gov": webbook.ScrapeWebbook{}}
var casFormat, _ = regexp.Compile("[0-9]{2,7}-[0-9]{2,2}-[0-9]")

func main() {

	var input string
	var output string
	var source string
	var quick string

	flag.StringVar(&source, "source", "webbook.nist.gov", "Source sit for retrieving data, default chemsynt")
	flag.StringVar(&input, "input", "", "Source CSV file for input. Expected to have single column or column label of 'Cas number'\nPossible values:")
	flag.StringVar(&output, "output", "", "Output CSV file containing results of parsing. Will contain original info from CSV")
	flag.StringVar(&quick, "quick", "", "Provide single CAS id for testing or quick reference")

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
	headerRow := checkHeaderRow(lines)
	casNumberColumn := checkCasColumn(lines, headerRow)
	//TODO: Refactor to call once for each CAS and allow cycling providers.
	outputLines, _ := availableScrapers[source].RunScrape(lines, headerRow, casNumberColumn)
	writeFile(outputLines, output)
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

func checkCasColumn(lines [][]string, headerRow bool) int {
	row := 0
	if headerRow {
		row = 1
	}
	for i, column := range lines[row] {
		if casFormat.MatchString(column) {
			fmt.Printf("Found CAS number column number on column %d.\n", i)
			return i
		}
	}
	fmt.Printf("No CAS-number found on row %d, do you have empty rows or rows not containing CAS number?\n", row)
	os.Exit(1)
	return 0
}

func checkHeaderRow(lines [][]string) bool {
	for _, column := range lines[0] {
		if casFormat.MatchString(column) {
			fmt.Println("Found CAS numbers on first row, assuming no header row.\n\r")
			return false
		}
	}
	fmt.Println("Found no CAS numbers on first row, assuming header row.\n\r")
	return true
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
		w.WriteAll(lines)
	}

}
