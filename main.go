package Endpoint

import (
	"fmt"
	"github.com/mitchelllisle/avista-ingest-flights/src/aviationStack"
	"github.com/mitchelllisle/avista-ingest-flights/src/config"
	"github.com/mitchelllisle/avista-ingest-flights/src/gcp"
	"github.com/mitchelllisle/avista-ingest-flights/src/utils"
	"log"
	"net/http"
)

func getSydneyArrivals(avista *aviationStack.AviationStack, storage *gcp.GCS, folder string) {
	var sydneyChannel = make(chan aviationStack.Flight)
	var params = make(map[string]string)
	params["arr_iata"] = "SYD"

	go avista.StreamFlights(sydneyChannel, params)

	payload := aviationStack.CollectRecords(sydneyChannel)

	folder = fmt.Sprintf("%s/sydney", folder)
	fileName := utils.GenerateFileName(folder)

	log.Println(fmt.Sprintf("writing file %s", fileName))
	storage.UploadFile(payload, fileName, "application/json")

}

func getLondonDepartures(avista *aviationStack.AviationStack, storage *gcp.GCS, folder string) {
	var londonChannel = make(chan aviationStack.Flight)
	var params = make(map[string]string)
	params["dep_iata"] = "LHR"

	go avista.StreamFlights(londonChannel, params)

	payload := aviationStack.CollectRecords(londonChannel)

	folder = fmt.Sprintf("%s/london", folder)
	fileName := utils.GenerateFileName(folder)

	log.Println(fmt.Sprintf("writing file %s", fileName))
	storage.UploadFile(payload, fileName, "application/json")
}

func main() {
	AppConfig := config.NewAppConfig()

	avista := aviationStack.InitAvista(
		AppConfig.AvistaConfig.URL,
		AppConfig.AvistaConfig.AccessToken,
		AppConfig.AvistaConfig.PageLimit)

	storage := gcp.InitGCS(AppConfig.GCPConfig.Project, AppConfig.GCPConfig.StorageBucket)

	getSydneyArrivals(avista, storage, AppConfig.GCPConfig.StorageFolder)
	getLondonDepartures(avista, storage, AppConfig.GCPConfig.StorageFolder)

	log.Println("finished")
}


func AvistaIngestFlights(w http.ResponseWriter, r *http.Request) {
	main()
}
