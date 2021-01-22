package aviationStack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mitchelllisle/avista-ingest-flights/src/utils"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"sync"
)

type AviationStack struct {
	URL string
	AccessToken string
	PageLimit int
}

func InitAvista(url, accessToken string, pageLimit int) *AviationStack {
	return &AviationStack{
		URL:         url,
		AccessToken: accessToken,
		PageLimit:   pageLimit,
	}
}

func (a *AviationStack) getPage(offset int, arrivalPortCode string) AviationStackFlightsResponse {
	url := fmt.Sprintf(
		"%s/flights?access_key=%s&limit=%d&offset=%d&arr_iata=%s",
		a.URL,
		a.AccessToken,
		a.PageLimit,
		offset,
		arrivalPortCode)

	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var flightResponse AviationStackFlightsResponse
	err := json.Unmarshal(body, &flightResponse)
	utils.PanicOnError(err, "Unable to unmarshall AviationStackFlightsResponse")
	return flightResponse
}

func getMaxPages(limit int, total int) int {
	if total > limit {
		totalPages := math.Ceil(float64(total) / float64(limit))
		return int(totalPages)
	} else {
		return 0
	}
}

func (a *AviationStack) getPageAndStreamRecords(offset int, arrivalPortCode string, channel chan <- Flight, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println(fmt.Sprintf("at offset: %d", offset))
	flights := a.getPage(offset, arrivalPortCode)
	for _, flight := range flights.Data {
		channel <- flight
	}
}


func (a *AviationStack) StreamFlights(channel chan <- Flight, arrivalPortCode string) {
	defer close(channel)
	var wg sync.WaitGroup

	flights := a.getPage(0, arrivalPortCode)
	for _, flight := range flights.Data {
		channel <- flight
	}

	totalPages := getMaxPages(flights.Pagination.Limit, flights.Pagination.Total)
	wg.Add(totalPages)

	for i := 1; i <= totalPages; i++ {
		go a.getPageAndStreamRecords(i * 100, arrivalPortCode, channel, &wg)
	}

	wg.Wait()
}

// Collect all records from channel and convert to []byte which is a json lines file when written
func CollectRecords(channel chan Flight) []byte {
	buffer := new(bytes.Buffer)

	for record := range channel {
		err := json.NewEncoder(buffer).Encode(record)
		utils.ContinueOnError(err, "Unable to marshal Flights into byte slice. Row will be discarded")
	}
	return buffer.Bytes()
}
