package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type payload struct {
	Results data
}

type data struct {
	Servant []servant      `json:"servants"`
	CE      []craftEssence `json:"ce"`
}

// type craftEssenceCollection []craftEssence

//USE REQ FORM PARSE FORM
//marshall struct data to json
//run this go server as an individual backend separate from node/react

//make one route that fetches all character names
//go should handle most of the scraping

func main() {
	r := mux.NewRouter()
	// routes consist of a path and a handler function.
	r.HandleFunc("/testGoApi", fetchOverallScrape).Methods("GET")
	r.HandleFunc("/fetchCard", fetchCard).Methods("GET")

	// bind to a port and pass our router in
	http.Handle("/", &middleWareServer{r})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func fetchCard(w http.ResponseWriter, r *http.Request) {
	if r != nil {
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}

		query := r.Form
		// fmt.Printf("type of form %T", query)
		// fmt.Println("")
		// fmt.Printf("value of form %v", query)
		// fmt.Println("")
		fmt.Println("card type", query["cardType"])
		fmt.Println("url value", query["pageLink"])

		individualScrapeMain(query["cardType"], query["pageLink"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		temp := []byte("this is fetchCard")
		w.Write(temp)
	}
}

func fetchOverallScrape(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	rTime := new(big.Int)
	fmt.Println(rTime.Binomial(1000, 10))

	w.Header().Set("Access-Control-Allow-Origin", "*")

	// fmt.Println("this is value of servantsMain", s)
	s := servantsMain()
	c := ceMain()

	d := data{Servant: s, CE: c}
	p := payload{d}

	searchJSON, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("this is value of searchJson", p)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(searchJSON)

	elapsed := time.Since(start)
	log.Printf("craft essence search took %s", elapsed)
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
