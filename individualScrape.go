package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type servantCard struct {
	EngName   string              `json:"engName"`
	JpName    string              `json:"jpName"`
	PageLink  string              `json:"pageLink"`
	Icon      string              `json:"icon"`
	Rarity    string              `json:"rarity"`
	CardType  string              `json:"cardType"`
	Available servantAvailability `json:"availability"`
	Params    servantParameters   `json:"params"`
}

type servantAvailability struct {
	jpContent bool
	event     bool
}

type servantParameters struct {
	minAtttack     int
	maxAtttack     int
	grailAttack    int
	minHP          int
	maxHP          int
	grailHP        int
	cost           int
	starAbsorption int
	starGen        float64
	npChargeAtk    float64
	npChargeDef    float64
	deathRate      float64
	growthCurve    string
	attribute      string
	alignment      string
}

type servantStats struct {
}

func main() {
	// determine logic here to create card type of servant or craft essence
	// will need two structs

	servant := servantCard{
		Available: servantAvailability{},
		Params:    servantParameters{},
	}

	scrapeServantPage(&servant)

	fmt.Println(servant)
}

func scrapeServantPage(s *servantCard) {
	// https://grandorder.wiki/Jeanne_d%27Arc#Official
	// https://grandorder.wiki/Ryougi_Shiki_(Assassin)
	// https://grandorder.wiki/Elisabeth_B%C3%A1thory_(Halloween)
	doc, err := goquery.NewDocument("https://grandorder.wiki/Ryougi_Shiki_(Assassin)")
	if err != nil {
		log.Fatal(err)
	}

	// separate the doc.Finds into go routines if they differ by search, like search tbody, tr, etc...
	// use a pointer to a struct to set the values
	// use a waitgroup of 1 to wait for all of the struct values to be set on the pointer

	//the first two tbodys will have the jp only content and the event reward servant tags...if they exist

	var wg sync.WaitGroup

	// create the length of waitgroups depending on whether it's a card or not
	// or just use channels

	wg.Add(2)
	go servantTags(doc, &wg, s)
	go servantParams(doc, &wg, s)

	wg.Wait()

	//put all of this in a channel i think
	//split up the tasks concurrently...would be good
	//one function handle alerts
	//one function handle skills
	//one function handle stats etc
}

func servantTags(doc *goquery.Document, wg *sync.WaitGroup, s *servantCard) {
	doc.Find("tbody").Each(func(index int, item *goquery.Selection) {
		// use index to determine if jp event only/event character only

		if index == 1 || index == 2 {
			// bTagExist := false
			testCaseJP := "JP-only Content"
			testCaseEvent := "Event Reward Servant"

			query := item.Find("b").First().Text()

			switch query {
			case testCaseJP:
				s.Available.jpContent = true
				fmt.Println("is jpOnly")
			case testCaseEvent:
				s.Available.event = true
				fmt.Println("is event servant")
			}
		}

	})
	defer wg.Done()
}

func servantParams(doc *goquery.Document, wg *sync.WaitGroup, s *servantCard) {
	doc.Find("div[style*='width:300px;margin-top:10px']").First().Find("div[style*='display:inline-block;width:150px;text-align:left;padding: 3px 0 3px 6px;border-left: 2px #4b9acc solid;vertical-align:top']").Each(func(index int, item *goquery.Selection) {
		// use index to determine if jp event only/event character only

		switch index {
		case 0:
			min, max, grail := processParams(item.Text())
			// move all this to the procesParams func
			s.Params.minAtttack = min
			s.Params.maxAtttack = max
			s.Params.grailAttack = grail
			fmt.Println("attack")
		case 1:
			min, max, grail := processParams(item.Text())
			s.Params.minHP = min
			s.Params.maxHP = max
			s.Params.grailHP = grail
			fmt.Println("hp")
		case 2:
			i, err := strconv.ParseFloat(item.Text(), 64)
			if err != nil {
				panic(err)
			}
			s.Params.cost = int(i)
		case 3:
			i, err := strconv.ParseFloat(item.Text(), 64)
			if err != nil {
				panic(err)
			}
			s.Params.starAbsorption = int(i)
		case 4:
			temp := removeSymbols(item.Text())
			i, err := strconv.ParseFloat(temp, 64)
			if err != nil {
				panic(err)
			}
			s.Params.starGen = i
		case 5:
			temp := removeSymbols(item.Text())
			i, err := strconv.ParseFloat(temp, 64)
			if err != nil {
				panic(err)
			}
			s.Params.npChargeAtk = i
		case 6:
			temp := removeSymbols(item.Text())
			i, err := strconv.ParseFloat(temp, 64)
			if err != nil {
				panic(err)
			}
			s.Params.npChargeDef = i
		case 7:
			temp := removeSymbols(item.Text())
			i, err := strconv.ParseFloat(temp, 64)
			if err != nil {
				panic(err)
			}
			s.Params.deathRate = i
		case 8:
			s.Params.growthCurve = item.Text()
		case 9:
			s.Params.attribute = item.Text()
		case 10:
			s.Params.alignment = item.Text()
		}

		// fmt.Println(temp)
	})
	defer wg.Done()
}

func processParams(s string) (int, int, int) {

	parts := strings.Fields(s)
	tempMin := parseStrings(parts[0])
	tempMax := parseStrings(parts[2])

	var grailRegex = regexp.MustCompile(`\((.*?)\)`)
	grail := grailRegex.FindStringSubmatch(s)[1]
	tempGrail := parseStrings(grail)

	return tempMin, tempMax, tempGrail
}

func parseStrings(s string) int {
	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	temp := int(i)
	return temp
}

func removeSymbols(s string) string {
	regexp, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		panic(err)
	}
	processString := regexp.ReplaceAllString(s, "")
	return processString
}

func servantStats(doc *goquery.Document, wg *sync.WaitGroup, s *servantCard) {
	doc.Find("div[style*='width:300px;margin-top:10px']").First().Next().Find("div[style*='display:inline-block;width:150px;text-align:left;padding: 3px 0 3px 6px;border-left: 2px #4b9acc solid;vertical-align:top']").Each(func(index int, item *goquery.Selection) {

		fmt.Println(item.Text())
	})
	defer wg.Done()
}
