package main

import (
	"encoding/csv"
	"github.com/gocolly/colly"
	"log"
	"os"
)

func main(){

	fName := "cryptocoinmarketcap.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{"Name"})

	// Instantiate default collector
	c := colly.NewCollector()

	c.OnHTML(".cmc-table__table-wrapper-outer", func(e *colly.HTMLElement) {
		writer.Write([]string{
			e.ChildText("div.cmc-table.column-name.cmc-table.column-name--narrow-layout.sc-1kxikfi-0.eTVhdN a"),
		})
	})

	c.Visit("https://coinmarketcap.com/all/views/all/")

	log.Printf("Scraping finished, check file %q for results\n", fName)
}

