package main

import (
	"finplatform/finjson"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"time"
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

	finJson := &finjson.FinJson{ApiKey: apiKey}

	data, err := finJson.ParseJson("OPEC", "ORB", "2010-01-01", "2020-01-01")
	if err != nil {
		panic(err)
	}

	meta := data.Dataset.MetaData
	fmt.Printf("Data Summary:\nTicker: %s\nDataset: %s\nStart Date: %s\nEnd Date: %s\n",
	 meta.DatasetCode, meta.DatabaseCode,
	  MakeDateString(data.Dataset.StartDate.Date()),
	   MakeDateString(data.Dataset.EndDate.Date()))
}
