package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var PORT = ":8080"

type Cuaca struct {
	Status struct {
		Water int `json:"water"`
		Wind  int `json:"wind"`
	}
}

type Kondisi struct {
	AirCond   string
	WaterCond string
}

var info = []byte(`
		{
			"status":{
				"water": 5,
				"wind": 10
			}
		}
	`)

func repetition(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

var kondisi Kondisi

func updateData(t time.Time) {

	var cuaca Cuaca

	err := json.Unmarshal(info, &cuaca)
	if err != nil {
		log.Fatalln("error:", err)
	}

	cuaca.Status.Water = rand.Intn(10)
	cuaca.Status.Wind = rand.Intn(10)

	_, err = json.Marshal(&cuaca)
	if err != nil {
		panic(err)
	}

	if cuaca.Status.Water < 5 {
		kondisi.WaterCond = "Aman"
		// fmt.Println("Status Air: Aman")
	} else if cuaca.Status.Wind < 6 {
		kondisi.AirCond = "Aman"
		// fmt.Println("Status Udara: Aman")
	} else if cuaca.Status.Water > 6 && cuaca.Status.Water < 8 {
		kondisi.WaterCond = "Siaga"
		// fmt.Println("Status Air: Siaga")
	} else if cuaca.Status.Wind > 7 && cuaca.Status.Wind < 15 {
		kondisi.AirCond = "Siaga"
		// fmt.Println("Status Angin: Siaga")
	} else {
		kondisi.WaterCond = "bahaya"
		kondisi.AirCond = "bahaya"
		// fmt.Println("Status: Bahaya")
	}
	// fmt.Printf("%s", kondisi)
}

func htmlData(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tpl, err := template.ParseFiles("template.html")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		tpl.Execute(w, kondisi)
		return

	}
	http.Error(w, "Invalid Method", http.StatusBadRequest)
}

func main() {

	go repetition(15*time.Second, updateData)
	http.HandleFunc("/data/html", htmlData)

	fmt.Println("Application is listening on port ", PORT)
	http.ListenAndServe(PORT, nil)

}
