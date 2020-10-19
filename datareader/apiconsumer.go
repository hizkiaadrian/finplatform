package datareader

type ApiConsumer interface {
	SetApiKey(apiKey string)

	FetchData(databaseName string, ticker string, startDate string, endDate string) (Data, error)
}

