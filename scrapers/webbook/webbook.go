package webbook

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/antchfx/htmlquery"
	"github.com/parnurzeal/gorequest"
)

const url = "https://webbook.nist.gov/cgi/cbook.cgi?ID=%s&Units=SI" //Cas number inserted in ID
var request = gorequest.New()

type ScrapeWebbook struct {
}

type Result struct {
	Name  string
	Value string
}

/*
*Attempts to parse data on CAS numbers provided in CSV format.
 */
func (s ScrapeWebbook) RunScrape(inputLines [][]string, headerRow bool, casColumn int) ([][]string, error) {
	startLine := 0
	if headerRow {
		startLine = 1
	}
	results := make([][]Result, len(inputLines)-startLine)
	fmt.Printf("Runing webbook scraper %d %s.\n", len(inputLines), inputLines[0:0])
	for i, line := range inputLines[startLine:] {
		result, err := parseWebsite(line[casColumn])
		if err != nil || results == nil || len(results) == 0 {
			fmt.Printf("Could not retrieve data for CAS-nr: %s: %s", line[casColumn], err.Error())
		}
		results[i] = result

	}
	return createOutputLines(inputLines, headerRow, results), nil

}

func parseWebsite(casNumber string) (results []Result, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic occured: ", r)
			err = errors.New("Panic when parsing website")
		}
	}()
	targetURL := fmt.Sprintf(url, casNumber)
	fmt.Printf("Retrieving data for %s\n", targetURL)

	resp, errs := http.Get(targetURL)
	if errs != nil {
		return nil, errs
	}
	results, err = parseBody(resp)
	return results, err
}

func parseBody(response *http.Response) ([]Result, error) {
	fmt.Println("Parsing webbook result..")
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return matchXpath(string(bytes))
}

func matchXpath(body string) ([]Result, error) {

	doc, err := htmlquery.Parse(strings.NewReader(body))
	if err != nil {
		return []Result{}, err
	}

	main := htmlquery.FindOne(doc, "//main")
	if main == nil {
		return []Result{Result{Name: "Molecular Weight", Value: "Could not parse"}}, nil
	}
	molecularWeight := htmlquery.FindOne(main, "/ul[1]/li[2]").LastChild
	checkHeader := htmlquery.FindOne(main, "/ul[1]/li[2]/strong[1]/a[1]").LastChild

	fmt.Printf("Mol: %+v\n", molecularWeight)
	fmt.Printf("header: %+v\n", checkHeader)

	return []Result{Result{Name: checkHeader.Data, Value: molecularWeight.Data}}, nil

}
func createOutputLines(inputLines [][]string, headerRow bool, results [][]Result) [][]string {

	header := 0
	if headerRow {
		header = 1
	}

	//Useful for headers
	line := make([]string, len(inputLines[0]))
	if headerRow {
		line = inputLines[0]
	}

	lines := make([][]string, len(inputLines)+1)

	for _, res := range results[0] {
		line = append(line, res.Name)
	}

	lines[0] = line
	for i, lineResults := range results {
		newLine := inputLines[i+header]
		for _, result := range lineResults {
			newLine = append(newLine, strings.TrimSpace(result.Value))
		}
		lines[i+1] = newLine
	}
	fmt.Printf("Lines %s", lines)
	return lines
}
