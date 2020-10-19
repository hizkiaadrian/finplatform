package datareader

import (
	"encoding/json"
	"errors"
	"finplatform/dateutils"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type JsonApiConsumer struct {
	_key string
}

func (jac *JsonApiConsumer) SetApiKey(apiKey string) {
	(*jac)._key = apiKey
}

func (jac *JsonApiConsumer) FetchData(datasetName string, ticker string, startDate string, endDate string) (Data, error) {
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

	q := req.URL.Query()
	q.Add("start_date", startDate)
	q.Add("end_date", endDate)
	q.Add("api_key", (*jac)._key)
	req.URL.RawQuery = q.Encode()

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

	var d JsonData
	err = json.Unmarshal(body, &d)
	if err != nil {
		return nil, err
	}

	if len(d.Dataset.DatasetCode) == 0 {
		return nil, errors.New("Data not available. Please ensure that ticker symbol is valid, and that the dates are in YYYY-MM-DD")
	}

	return d.Dataset, nil
}

type JsonData struct {
	Dataset Dataset `json:"dataset"`
}

type Dataset struct {
	MetaData
	StartDate dateutils.DateOnly `json:"start_date"`
	EndDate   dateutils.DateOnly `json:"end_date"`
	PriceData []PriceData        `json:"data"`
}

type MetaData struct {
	DatasetCode         string             `json:"dataset_code"`
	DatabaseCode        string             `json:"database_code"`
	Name                string             `json:"name"`
	Description         string             `json:"description"`
	Frequency           string             `json:"frequency"`
	Type                string             `json:"type"`
	RefreshedAt         time.Time          `json:"refreshed_at"`
	NewestAvailableDate dateutils.DateOnly `json:"newest_available_date"`
	OldestAvailableDate dateutils.DateOnly `json:"oldest_available_date"`
}

type PriceData struct {
	Date  dateutils.DateOnly
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
