package aviationStack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mitchelllisle/avista-ingest-flights/src/utils"
	"github.com/wlredeye/jsonlines"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
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

func (a *AviationStack) getPage(offset int, params map[string]string) AviationStackFlightsResponse {
	url := fmt.Sprintf("%s/flights", a.URL)

	req, _ := http.NewRequest("GET", url, nil)
	queryParams := req.URL.Query()
	queryParams.Add("access_key", a.AccessToken)
	queryParams.Add("offset", strconv.Itoa(offset))
	queryParams.Add("limit", strconv.Itoa(a.PageLimit))
	for key, value := range params {
		queryParams.Add(key, value)
	}
	req.URL.RawQuery = queryParams.Encode()

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

func (a *AviationStack) getPageAndStreamRecords(offset int, params map[string]string, channel chan <- Flight, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println(fmt.Sprintf("at offset: %d", offset))
	flights := a.getPage(offset, params)
	for _, flight := range flights.Data {
		channel <- flight
	}
}


func (a *AviationStack) StreamFlights(channel chan <- Flight, params map[string]string) {
	defer close(channel)
	var wg sync.WaitGroup

	flights := a.getPage(0, params)
	for _, flight := range flights.Data {
		channel <- flight
	}

	totalPages := getMaxPages(flights.Pagination.Limit, flights.Pagination.Total)
	wg.Add(totalPages)

	for i := 1; i <= totalPages; i++ {
		go a.getPageAndStreamRecords(i * 100, params, channel, &wg)
	}

	wg.Wait()
}

// Collect all records from channel and convert to []byte which is a json lines file when written
func CollectRecords(channel chan Flight) []byte {
	var records []Flight
	var buffer bytes.Buffer

	for record := range channel {
		records = append(records, record)
	}
	err := jsonlines.Encode(&buffer, &records)
	utils.PanicOnError(err, "unable to parse data into json lines")
	return buffer.Bytes()
}
