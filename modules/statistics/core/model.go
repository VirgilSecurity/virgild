package core

type Endpoint int

const (
	UnknowEndpoint      Endpoint = 0
	GetCardEndpoint     Endpoint = 1
	SearchCardsEndpoint Endpoint = 2
	CreateCardEndpoint  Endpoint = 3
	RevokeCardEndpoint  Endpoint = 4
)

type RequestStatistics struct {
	Id        int64  `json:"id"`
	Date      int64  `json:"date"`
	DateMonth int64  `json:"-"`
	Token     string `json:"token"`
	Method    string `json:"method"`
	Resource  string `json:"resource"`
}

type StatisticDayGroup int

const (
	Hour  StatisticDayGroup = 0
	Day   StatisticDayGroup = 1
	Month StatisticDayGroup = 2
)

type StatisticGroup struct {
	Count    int
	Date     int64
	Token    string
	Endpoint Endpoint
}
