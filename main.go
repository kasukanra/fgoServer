package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type payload struct {
	values data
}

type data struct {
	Servent servant      `json:"servantName"`
	CE      craftEssence `json:"ceName"`
}

type servant []string
type craftEssence []string

//USE REQ FORM PARSE FORM
//marshall struct data to json
//run this go server as an individual backend separate from node/react

//make one route that fetches all character names
//go should handle most of the scraping

func main() {
	r := mux.NewRouter()
	// routes consist of a path and a handler function.
	r.HandleFunc("/testGoApi", fetchOverallScrape).Methods("GET")

	// bind to a port and pass our router in
	http.Handle("/", &middleWareServer{r})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func fetchOverallScrape(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// fmt.Printf("type of request %T", r)
	// fmt.Println("")
	// fmt.Printf("value of request %v", r)
	// fmt.Println("")
	// fmt.Println("request body", r.Body)
	//check form parse req
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	query := r.Form
	fmt.Printf("type of form %T", query)
	fmt.Println("")
	fmt.Printf("value of form %v", query)
	fmt.Println("")
	fmt.Println("url value", query["name"])

	var servant []string
	var ce []string

	servant = append(servant, "merlin")
	ce = append(ce, "formalcraft")
	servant = append(servant, "ishtar")

	d := data{servant, ce}
	p := payload{d}
	searchJSON, err := json.Marshal(p)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(searchJSON)
}

type middleWareServer struct {
	r *mux.Router
}

func (s *middleWareServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}
	// Lets Gorilla work
	s.r.ServeHTTP(rw, req)
}
