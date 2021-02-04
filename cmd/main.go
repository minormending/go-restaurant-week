package main

import (
	"errors"
	"flag"
	"os"

	"github.com/minormending/go-restaurant-week/client"
	"github.com/minormending/go-restaurant-week/formatters"
)

var (
	overwrite = flag.Bool("overwrite", false, "overwrite the outfile if it already exists.")
	help      = flag.Bool("help", false, "prints help information")
)

var usage = `Usage: restaurant-week [OPTIONS] API_KEY OUTFILE

Generates an HTML map of the NYC restuarant week restaurants.

Arguments:
	API_KEY		Google Maps API V3 key
	OUTFILE		filename to save the HTML file

Options:
	-overwrite bool
		overwrites the outfile if it already exists.
	-help bool
		prints help information
`

func main() {
	flag.Parse()
	apiKey := flag.Arg(0)
	filename := flag.Arg(1)

	restaurants, err := client.GetRestaurantInfo()
	if err != nil {
		panic(err)
	}

	info, err := os.Stat(filename)
	if os.IsNotExist(err) == false || info.IsDir() {
		if *overwrite {
			if err = os.Remove(filename); err != nil {
				panic(err)
			}
		} else {
			panic(errors.New("outfile file already exists and overwrite not specified"))
		}
	}

	file, err := os.OpenFile(filename, os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0666)
	defer file.Close()

	if err = formatters.ToHTML(file, apiKey, restaurants); err != nil {
		panic(err)
	}
}
