package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
)

type cityTable struct {
	Starting_point string `json:"starting_point"`
	Destination    string `json:"destination"`
	Route          string `json:"route"`
	Bus_num        string `json:"bus_num"`
	Time           string `json:"time"`
}

type intercityTable struct {
	Starting_point string `json:"starting_point"`
	Destination    string `json:"destination"`
	Starting_time  string `json:"starting_time"`
	Arrival_time   string `json:"arrival_time"`
	Price          string `json:"price"`
	Using_time     string `json:"using_time"`
}

var json_value any

func jsonContentTypeMiddleware(next http.Handler, classification string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := r.URL.Query().Get("start")
		destination := r.URL.Query().Get("destination")

		fmt.Println(start, destination)

		fetchData(start, destination, classification)
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func fetchData(start string, destination string, classification string) {
	dsn := "root:1234@tcp(127.0.0.1:3306)/Bus_TimeTable?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	var cityTimeTable []cityTable
	var intercityTimeTable []intercityTable

	if err != nil {
		panic("Db 연결에 실패하였습니다.")
	}

	if classification == "city" {
		db.Table("cityBusTable").Where("starting_point like ? and destination like ?", "%"+start+"%", "%"+destination+"%").Find(&cityTimeTable)
		json_value = cityTimeTable
	} else if classification == "limit" {
		db.Table("cityBusTable").Where("starting_point like ? and destination like ? and time > curtime()", "%"+start+"%", "%"+destination+"%").Order("time asc").Limit(1).Find(&cityTimeTable)
		json_value = cityTimeTable
	} else if classification == "intercity" {
		db.Table("intercityBusTable").Where("starting_point like ? and destination like ?", "%"+start+"%", "%"+destination+"%").Find(&intercityTimeTable)
		json_value = intercityTimeTable
	}
}

func main() {
	mux := http.NewServeMux()
	userHandler := http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		wr.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		wr.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		json.NewEncoder(wr).Encode(json_value)
	})

	mux.Handle("/city", jsonContentTypeMiddleware(userHandler, "city"))
	mux.Handle("/intercity", jsonContentTypeMiddleware(userHandler, "intercity"))
	mux.Handle("/limit", jsonContentTypeMiddleware(userHandler, "limit"))
	http.ListenAndServe(":8889", mux)
}
