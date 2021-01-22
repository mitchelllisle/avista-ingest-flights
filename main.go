package main

import (
	"fmt"
	"github.com/mitchelllisle/avista-ingest-flights/src/aviationStack"
	"github.com/mitchelllisle/avista-ingest-flights/src/config"
	"github.com/mitchelllisle/avista-ingest-flights/src/gcp"
	"log"
	"net/http"
	"time"
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

	payload := aviationStack.CollectRecords(flightChannel)

	fileName := fmt.Sprintf("%s.jsonl", time.Now().Format(time.RFC3339))

	log.Println(fmt.Sprintf("writing file %s", fileName))
	storage.UploadFile(payload, fileName, "application/json")

	log.Println("finished")
}


func AvistaIngestFlights(w http.ResponseWriter, r *http.Request) {
	main()
}
