package main

import (
	"finplatform/datareader"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"time"
	"strings"
)

func MakeDateString(year int, month time.Month, day int) string {
	return fmt.Sprintf("%d %s %d", day, month.String(), year)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	apiKey := os.Getenv("QUANDL_API_KEY")

	var apiConsumer datareader.ApiConsumer
	apiConsumer = new(datareader.JsonApiConsumer)
	apiConsumer.SetApiKey(apiKey)

	d, err := apiConsumer.FetchData("OPEC", "ORB", "2010-01-01", "2020-01-01")
	if err != nil {
		panic(err)
	}
	data := d.(datareader.Dataset)

	apiConsumer = new(datareader.CsvApiConsumer)
	apiConsumer.SetApiKey(apiKey)
	c, err := apiConsumer.FetchData("WIKI", "AAPL", "2010-01-01", "2020-01-01")
	if err != nil {
		panic(err)
	}
	csv := c.([]string)

	for _, s := range strings.Split(csv[0],`,`) {
		fmt.Println(s)
	}

	meta := data.MetaData
	fmt.Printf("Data Summary:\nTicker: %s\nDataset: %s\nStart Date: %s\nEnd Date: %s\n",
	 meta.DatasetCode, meta.DatabaseCode,
	  MakeDateString(data.StartDate.Date()),
	   MakeDateString(data.EndDate.Date()))
}
