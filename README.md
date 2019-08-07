# CAS number scraper

Project to scrape molecular weights for a provided CSV of CAS-Numbers.
This is helpful to friends of the programmer but open to anyone else with such a need.

## Building the scraper

- Clone the git repository
- If cloned outside of GOPATH add the directory to the GOPATH, for example using export GOPATH=$GOPATH;$(pwd)
- Build the binary using `go build`

## Using the scraper

In the directory with the binary:

- `go --help` shows help prompt.
To run a run using the provided testfiles:
- `go -input testfiles/test.csv -output output.csv` this will print the results to output.csv.
- `go -input testfiles/test.csv -output -` this will print to stdout.
- `go -quick 64-17-5` quick check for single CAS number. Defaults to printing to stdout.
output

To check a simple cas number 




