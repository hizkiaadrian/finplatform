package main

import (
	"finplatform/finjson"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	apiKey := os.Getenv("QUANDL_API_KEY")

	finJson := &finjson.FinJson{ApiKey: apiKey}

	test, err := finJson.ParseJson("OPEC", "ORB", "", "")
	if err != nil {
		panic(err)
	}

	i := test.Dataset.MetaData
	fmt.Printf("Data for ticker %s from %s dataset loaded succesfully", i.DatasetCode, i.DatabaseCode)
}
