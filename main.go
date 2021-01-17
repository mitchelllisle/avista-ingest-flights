package main

import (
	"avista-ingest-flights/src/aviationStack"
	"avista-ingest-flights/src/config"
	"avista-ingest-flights/src/gcp"
	"avista-ingest-flights/src/utils"
	"encoding/json"
)

func main() {
	AppConfig := config.NewAppConfig()

	avista := aviationStack.InitAvista(
		AppConfig.AvistaConfig.URL,
		AppConfig.AvistaConfig.AccessToken,
		AppConfig.AvistaConfig.PageLimit)

	storage := gcp.InitGCS(AppConfig.GCPConfig.Project, AppConfig.GCPConfig.StorageBucket)

	var flightChannel = make(chan aviationStack.Flight)

	go avista.StreamFlights(flightChannel, AppConfig.AvistaConfig.ArrivalCode)

	var records []aviationStack.Flight

	for record := range flightChannel {
		records = append(records, record)
	}
	payload, err := json.Marshal(records)
	utils.PanicOnError(err, "Unable to marshal Flights into byte slice")

	storage.UploadFile(payload, "test", "application/json")
}
