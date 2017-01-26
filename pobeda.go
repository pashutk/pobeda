package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	grequests "github.com/levigross/grequests"
	html "golang.org/x/net/html"
)

type SessionOptions struct {
	MarketType  string
	FromStation string
	ToStation   string
	BeginDate   string
}

type Flight struct {
	Date  string
	Price int
}

func (flight Flight) String() string {
	return fmt.Sprintf("Flight date: %s, price: %d", flight.Date, flight.Price)
}

type ByPrice []Flight

func (s ByPrice) Len() int { return len(s) }
func (s ByPrice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByPrice) Less(i, j int) bool {
	return s[i].Price < s[j].Price
}

func getURLForSession(options SessionOptions) string {
	const urlTemplate = "https://booking.pobeda.aero/ExternalSearch.aspx?marketType=%s&fromStation=%s&toStation=%s&beginDate=%s&adultCount=1&childrenCount=0&infantCount=0&currencyCode=RUB&utm_source=pobeda&culture=ru-RU"

	return fmt.Sprintf(urlTemplate, options.MarketType, options.FromStation, options.ToStation, options.BeginDate)
}

func initSession(session grequests.Session, options SessionOptions) string {
	url := getURLForSession(options)
	resp, err := session.Get(url, nil)

	if err != nil {
		log.Fatalln("Unable to make request: ", err)
	}

	return resp.String()
}

func getMonthPricesHTML(session grequests.Session, month time.Time) string {
	selectedDate := month.Format("2006-01-02")
	resp, err := session.Post("https://booking.pobeda.aero/AjaxMonthLowFareAvailaibility.aspx", &grequests.RequestOptions{
		Data: map[string]string{"dateSelected": selectedDate},
	})

	if err != nil {
		log.Fatalln("Unable to make request: ", err)
	}

	return resp.String()
}

func isDayMonthDiv(token html.Token) bool {
	result := false
	for _, attr := range token.Attr {
		if attr.Key == "data-type" && attr.Val == "dayMonth" {
			result = true
		}
	}
	return result
}

func isDayMonthWithFlightsDiv(token html.Token) bool {
	if isDayMonthDiv(token) {
		result := false
		for _, attr := range token.Attr {
			if attr.Key == "data-hasflights" && attr.Val == "true" {
				result = true
			}
		}
		return result
	}
	return false
}

func isPriceDiv(token html.Token) bool {
	result := false
	for _, attr := range token.Attr {
		if attr.Key == "class" && attr.Val == "price" {
			result = true
		}
	}
	return result
}

func getDataDateAttr(token html.Token) string {
	result := ""
	for _, attr := range token.Attr {
		if attr.Key == "data-date" {
			result = attr.Val
		}
	}
	return result
}

func getAttributeAttr(token html.Token) string {
	result := ""
	for _, attr := range token.Attr {
		if attr.Key == "attribute" {
			result = attr.Val
		}
	}
	return result
}

func parsePriceFromDataAttr(token html.Token) int {
	str := getAttributeAttr(token)
	result, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Fatalln("Unable to parse price: ", err)
	}
	return int(result)
}

func parsePrices(monthPrices string) []Flight {
	reader := strings.NewReader(monthPrices)
	tokenizer := html.NewTokenizer(reader)

	creatingFlight := false

	var flights []Flight
	flight := Flight{}

	for {
		tt := tokenizer.Next()

		switch {
		case tt == html.ErrorToken:
			return flights
		case tt == html.StartTagToken:
			t := tokenizer.Token()

			if isDayMonthWithFlightsDiv(t) {
				creatingFlight = true
				flight.Date = getDataDateAttr(t)
			}

			if creatingFlight && isPriceDiv(t) {
				flight.Price = parsePriceFromDataAttr(t)
				flights = append(flights, flight)
				flight = Flight{}
				creatingFlight = false
			}
		}
	}
}

func removeDuplicatesFromFlights(flights []Flight) []Flight {
	uniqueSet := map[string]Flight{}
	for _, flight := range flights {
		if _, ok := uniqueSet[flight.Date]; ok {
			continue
		} else {
			uniqueSet[flight.Date] = flight
		}
	}
	var result []Flight
	for _, flight := range uniqueSet {
		result = append(result, flight)
	}
	return result
}

func main() {
	fmt.Println("Poehali!")

	session := grequests.NewSession(nil)

	from := "VKO"
	to := "LCA"
	if len(os.Args) > 1 {
		to = os.Args[1]
	}
	currentTime := time.Now().Local()
	beginDate := currentTime.Format("2006-01-02")

	options := SessionOptions{
		MarketType:  "OneWay",
		FromStation: from,
		ToStation:   to,
		BeginDate:   beginDate,
	}

	initSession(*session, options)

	var flights []Flight
	MONTH_COUNT := 6
	for i := 0; i < MONTH_COUNT; i++ {
		nextMonth := currentTime.AddDate(0, i, 0)
		html := getMonthPricesHTML(*session, nextMonth)
		flights = append(flights, parsePrices(html)...)
	}

	flights = removeDuplicatesFromFlights(flights)
	sort.Sort(ByPrice(flights))

	for _, flight := range flights {
		fmt.Println(flight)
	}

	fmt.Println("Pobeda!")
}
