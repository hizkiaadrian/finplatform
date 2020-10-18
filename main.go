package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type dateOnly struct {
	time.Time
}

func (t *dateOnly) UnmarshalJSON(buf []byte) error {
	tt, err := time.Parse("2006-01-02", strings.Trim(string(buf), `"`))
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}

type Data struct {
	Dataset Dataset `json:"dataset"`
}

type Dataset struct {
	MetaData
	StartDate dateOnly    `json:"start_date"`
	EndDate   dateOnly    `json:"end_date"`
	PriceData []PriceData `json:"data"`
}

type MetaData struct {
	DatasetCode         string    `json:"dataset_code"`
	DatabaseCode        string    `json:"database_code"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	Frequency           string    `json:"frequency"`
	Type                string    `json:"type"`
	RefreshedAt         time.Time `json:"refreshed_at"`
	NewestAvailableDate dateOnly  `json:"newest_available_date"`
	OldestAvailableDate dateOnly  `json:"oldest_available_date"`
}

type PriceData struct {
	Date  dateOnly
	Price float64
}

func (p *PriceData) UnmarshalJSON(data []byte) error {
	var slicedData []interface{}
	if err := json.Unmarshal(data, &slicedData); err != nil {
		return err
	}

	p.Date.UnmarshalJSON([]byte(strings.Trim(slicedData[0].(string), `"`)))

	p.Price = slicedData[1].(float64)
	
	return nil

}

func init() {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}
}

func main() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Something went wrong in the main function. Please contact hizkiaadrians@gmail.com for fix")
		}
	}()

	var url string

	if quandlApiKey, exists := os.LookupEnv("QUANDL_API_KEY"); exists {
		url = fmt.Sprintf("https://www.quandl.com/api/v3/datasets/OPEC/ORB.json?api_key=%s", quandlApiKey)
	} else {
		panic("Cannot find Quandl API key")
	}

	quandlClient := http.Client{Timeout: time.Second * 2}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := quandlClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var datasets Data
	err = json.Unmarshal(body, &datasets)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", res.Status)
}
