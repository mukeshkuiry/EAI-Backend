package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
)

// Stock represents the structure of each stock
type Stocks struct {
	Symbol          string  `json:"symbol"`
	OpenPrice       float64 `json:"openPrice"`
	RefreshInterval int     `json:"refreshInterval"`
}

func LoadStockData(w http.ResponseWriter, r *http.Request) {
	// allow all cors requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var stocks []Stocks
	data, err := ioutil.ReadFile("data.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &stocks)
	if err != nil {
		log.Fatal(err)
	}
	go startStockPriceUpdater(stocks, 0)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stocks)
}

func FetchAndStoreStock() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("API_KEY")

	c := polygon.New(apiKey)

	// set params
	params := models.GetGroupedDailyAggsParams{
		Locale:     models.US,
		MarketType: models.Stocks,
		Date:       models.Date(time.Date(2023, 3, 8, 0, 0, 0, 0, time.Local)),
	}.WithAdjusted(true)

	// make request
	res, err := c.GetGroupedDailyAggs(context.Background(), params)
	if err != nil {
		log.Fatal(err)
	}
	modifiedData := res.Results[:20]

	var newData []Stocks
	for _, v := range modifiedData {
		newData = append(newData, Stocks{
			Symbol:          v.Ticker,
			OpenPrice:       float64(int(v.Open*10000)) / 10000,
			RefreshInterval: rand.Intn(5) + 1,
		})
	}

	data, err := json.Marshal(newData)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("data.json", data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func startStockPriceUpdater(stocks []Stocks, counter int) {
	for i := range stocks {
		if counter%stocks[i].RefreshInterval == 0 {
			// Generate a random number between -5 and 5
			randNum := rand.Intn(11) - 5
			// Calculate the percentage change
			percentageChange := float64(randNum) / 100
			// Apply the percentage change to the open price
			stocks[i].OpenPrice += stocks[i].OpenPrice * percentageChange
			// upto 4 decimal places
			stocks[i].OpenPrice = float64(int(stocks[i].OpenPrice*10000)) / 10000
		}
	}
	data, err := json.Marshal(stocks)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("data.json", data, 0644)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	if counter < 100000000000 { // Add a condition to stop the recursion after a certain number of iterations
		startStockPriceUpdater(stocks, counter+1)
	}
}
