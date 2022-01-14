package scrape

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Stock struct {
	Name  string
	Date  string
	Value string
}

func Scrape(names []string) {
	var baseURL string = "https://www.macrotrends.net/stocks/charts/"

	var stocks []Stock
	channel := make(chan []Stock)

	for i := 0; i < len(names); i++ {
		URL := baseURL + names[i] + "/american-airlines-group/operating-income"
		go getPage(names[i], URL, channel)
	}

	for i := 0; i < len(names); i++ {
		extractedStocks := <-channel
		stocks = append(stocks, extractedStocks...)
	}

	writeStocks(stocks)
	fmt.Println("Done, extracted", len(stocks))
}

func writeStocks(stocks []Stock) {
	channel := make(chan []string)

	file, err := os.Create("stocks.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"Name", "Date", "Value"}
	wErr := w.Write(headers)
	checkErr(wErr)

	for _, stock := range stocks {
		go writeStockDetail(stock, channel)
	}

	for i := 0; i < len(stocks); i++ {
		stockData := <-channel
		writeErr := w.Write(stockData)
		checkErr(writeErr)
	}
}

func writeStockDetail(stock Stock, c chan<- []string) {
	c <- []string{stock.Name, stock.Date, stock.Value}
}

func getPage(name, url string, mainChannel chan<- []Stock) {
	var stocks []Stock
	channel := make(chan Stock)

	res, err := http.Get(url)
	checkErr(err)
	checkCode(res)
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	count := doc.Find("#style-1 > div:nth-child(2) > table > tbody > tr").Length()
	fmt.Println(count)

	searchCards := doc.Find("#style-1 > div:nth-child(2) > table")
	searchCards.Each(func(i int, card *goquery.Selection) {
		for a := 1; a <= count; a++ {
			go extractStock(a, name, card, channel)
		}
	})

	for i := 0; i < count; i++ {
		stock := <-channel
		stocks = append(stocks, stock)
	}

	mainChannel <- stocks
}

func extractStock(i int, name string, card *goquery.Selection, c chan<- Stock) {
	num := strconv.Itoa(i)
	date := CleanString(card.Find("tr:nth-child(" + num + ")>td:nth-child(1)").Text())
	value := CleanString(card.Find("tr:nth-child(" + num + ")>td:nth-child(2)").Text())

	c <- Stock{
		Name:  name,
		Date:  date,
		Value: value,
	}
}

func checkErr(err error) {

	if err != nil {
		log.Fatal(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatal("Request failed with Status:", res.StatusCode)
	}
}

// CleanString cleans a string
func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}
