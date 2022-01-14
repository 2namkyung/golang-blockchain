package main

import (
	"learngo/practice/scrape"
)

func main() {
	list := []string{"AAPL", "MCRI"}
	scrape.Scrape(list)
}
