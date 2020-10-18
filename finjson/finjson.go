package finjson

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type DateOnly struct {
	time.Time
}

func (t *DateOnly) UnmarshalJSON(buf []byte) error {
	tt, err := time.Parse("2006-01-02", strings.Trim(string(buf), `"`))
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}

type JsonData struct {
	Dataset Dataset `json:"dataset"`
}

type Dataset struct {
	MetaData
	StartDate DateOnly    `json:"start_date"`
	EndDate   DateOnly    `json:"end_date"`
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
	NewestAvailableDate DateOnly  `json:"newest_available_date"`
	OldestAvailableDate DateOnly  `json:"oldest_available_date"`
}

type PriceData struct {
	Date  DateOnly
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

type FinJson struct {
	ApiKey string
}

func (fj *FinJson) ParseJson(datasetName string, ticker string, startDate string, endDate string) (*JsonData, error) {
	var (
		err error

		req *http.Request
		res *http.Response
	)

	url := fmt.Sprintf("https://www.quandl.com/api/v3/datasets/%s/%s.json", datasetName, ticker)

	httpClient := http.Client{Timeout: time.Second * 2}

	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.URL.Query().Add("start_date", startDate)
	req.URL.Query().Add("end_date", endDate)
	req.URL.Query().Add("api_key", (*fj).ApiKey)

	res, err = httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var datasets JsonData
	err = json.Unmarshal(body, &datasets)
	if err != nil {
		return nil, err
	}

	if len(datasets.Dataset.DatasetCode) == 0 {
		return nil, errors.New("Data not available. Please ensure that ticker symbol is valid, and that the dates are in YYYY-MM-DD")
	}

	return &datasets, nil
}
