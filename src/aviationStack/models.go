package aviationStack

import "time"

type Port struct {
	Airport         string    `json:"airport"`
	Timezone        string    `json:"timezone"`
	Iata            string    `json:"iata"`
	Icao            string    `json:"icao"`
	Terminal        string    `json:"terminal"`
	Gate            string    `json:"gate"`
	Baggage         string    `json:"baggage"`
	Delay           int       `json:"delay"`
	Scheduled       time.Time `json:"scheduled"`
	Estimated       time.Time `json:"estimated"`
	Actual          time.Time `json:"actual"`
	EstimatedRunway time.Time `json:"estimated_runway"`
	ActualRunway    time.Time `json:"actual_runway"`
}

type Airline struct {
	Name string `json:"name"`
	Iata string `json:"iata"`
	Icao string `json:"icao"`
}

type FlightMeta struct {
	Number     string `json:"number"`
	Iata       string `json:"iata"`
	Icao       string `json:"icao"`
}

type Flight struct {
	FlightDate   string `json:"flight_date"`
	FlightStatus string `json:"flight_status"`
	Departure    Port `json:"departure"`
	Arrival 	 Port `json:"arrival"`
	Airline 	 Airline `json:"airline"`
	FlightMeta   FlightMeta `json:"flight"`
}

type Pagination struct {
	Limit int
	Offset int
	Count int
	Total int
}

type AviationStackFlightsResponse struct {
	Pagination Pagination
	Data []Flight
}