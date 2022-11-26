package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Details struct {
	Amount    int    `json: "amount"`
	TimeStamp string `json: "timestamp`
}

var details = []Details{
	// {Amount: 100, TimeStamp: "2022-01-05 T15:04:05-0700"},
	// {Amount: 200, TimeStamp: "2022-01-10 T15:04:05-0700"},
}

type Response struct {
	Sum   int
	Avg   float64
	Max   int
	Min   int
	Count int
}

type LocationDetails struct {
	Location string `json: "location"`
}

var locDetails = []LocationDetails{}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/transaction", getTransactions).Methods("GET")
	router.HandleFunc("/transaction", postTransactions).Methods("POST")
	router.HandleFunc("/statistics", statistics).Methods("GET")
	router.HandleFunc("/transaction", deleteTransactions).Methods("DELETE")
	router.HandleFunc("/location", GetLocation).Methods("GET")
	router.HandleFunc("/location", AddLocation).Methods("POST")
	router.HandleFunc("/location", UpdateLocation).Methods("PUT")

	http.ListenAndServe(":9000", router)
}

func getTransactions(w http.ResponseWriter, r *http.Request) {

	detail := &details

	err1 := json.NewEncoder(w).Encode(detail)
	if err1 != nil {
		panic("error is due to encoder")
	}

}

//Create Transaction
func postTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error while read the request", err)
	}
	fmt.Println("reqBody:: ", string(reqBody))

	var transaction Details
	json.Unmarshal(reqBody, &transaction)

	t, err := time.Parse(time.RFC3339, transaction.TimeStamp)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t)
	loc, _ := time.LoadLocation("Asia/Kolkata")
	given := t.In(loc)
	g := given.Add(60 * time.Second)
	givenTime := g.Format(time.UnixDate)
	fmt.Println("GivenTime: ", givenTime)
	current := time.Now()
	currentTime := current.Format(time.UnixDate)
	sec := current.Add(-60 * time.Second)
	fmt.Println("Current Time: ", currentTime)

	if sec.Before(current) && given.Before(current) {
		w.WriteHeader(http.StatusCreated)
		details = append(details, transaction)
		json.NewEncoder(w).Encode(details)
	} else if given.Before(sec) {
		w.WriteHeader(http.StatusNoContent)
	} else if given.After(current) {
		w.WriteHeader(http.StatusUnprocessableEntity)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

//GET the Statistic of Transaction
func statistics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	var res Response

	sum, count := 0, 0
	max := details[0].Amount
	min := details[0].Amount
	for _, v := range details {
		sum += v.Amount
		if v.Amount > max {
			max = v.Amount
		} else if v.Amount < min {
			min = v.Amount
		}
		count++

	}
	avg := float64(sum) / float64(len(details))

	fmt.Println("Sum", sum)
	fmt.Println("AVG", avg)
	fmt.Println("MAX", max)
	fmt.Println("MIN", min)
	fmt.Println("Count", count)
	res = Response{
		Sum:   sum,
		Avg:   avg,
		Max:   max,
		Min:   min,
		Count: count,
	}
	// json.Marshal(x)
	json.NewEncoder(w).Encode(res)

}

//Delete Transaction
func deleteTransactions(w http.ResponseWriter, r *http.Request) {
	details = []Details{}
	w.WriteHeader(http.StatusNoContent)
}

//GET Location
func GetLocation(w http.ResponseWriter, r *http.Request) {

	l := &locDetails

	err1 := json.NewEncoder(w).Encode(l)
	if err1 != nil {
		panic("error is due to encoder")
	}

}

//POST Method
func AddLocation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	// locDetails := []LocationDetails{}
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error while reading:", err)
	}
	fmt.Println("reqBody:: ", string(reqBody))

	var loc LocationDetails
	err1 := json.Unmarshal(reqBody, &loc)
	if err1 != nil {
		log.Println("Error while unmarshal:", err1)
	}
	locDetails = append(locDetails, loc)
	json.NewEncoder(w).Encode(&locDetails)
	fmt.Println("locDetails", locDetails)
}

//UPDATE Method
func UpdateLocation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	// locDetails := []LocationDetails{}
	// loc, _ := (mux.Vars(r)["location"])
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error in reading: ", err)
	}
	var newLocation LocationDetails

	err1 := json.Unmarshal(reqBody, &newLocation)
	if err1 != nil {
		log.Println("Error in unmarshal the data", err1)
	}

	for i, l := range locDetails {
		l.Location = newLocation.Location
		locDetails = append(locDetails[:i], l)
		json.NewEncoder(w).Encode(l)
	}
}
