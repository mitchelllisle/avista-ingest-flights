package main

import (
	"avista-ingest-flights/src/aviationStack"
	"avista-ingest-flights/src/config"
	"fmt"
)


func main() {
	AppConfig := config.NewAppConfig()

	avista := aviationStack.AviationStack{
		URL:         AppConfig.AvistaConfig.URL,
		AccessToken: AppConfig.AvistaConfig.AccessToken,
		PageLimit: AppConfig.AvistaConfig.PageLimit,
	}

	var flightChannel = make(chan aviationStack.Flight)

	go avista.StreamFlights(flightChannel, AppConfig.AvistaConfig.ArrivalCode)

	//Temporary, this is where we'd write to database
	counter := 1
	for range flightChannel {
		fmt.Println(counter)
		counter ++
	}
}
