package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type servantCard struct {
	EngName     string              `json:"engName"`
	JpName      string              `json:"jpName"`
	PageLink    string              `json:"pageLink"`
	Icon        string              `json:"icon"`
	Rarity      string              `json:"rarity"`
	CardType    string              `json:"cardType"`
	Available   servantAvailability `json:"availability"`
	Params      servantParameters   `json:"params"`
	Stats       servantStats        `json:"stats"`
	Traits      servantTraits       `json:"traits"`
	CardChoices servantCardChoices  `json:"cardChoices"`
	CardHits    servantCardHits     `json:"cardHits"`
	Skills      servantSkills       `json:"skills"`
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

type servantTraits struct {
	traits string
}

type servantCardChoices struct {
	quick  int
	arts   int
	buster int
}

type servantCardHits struct {
	quickHits  int
	artsHits   int
	busterHits int
	extraHits  int
}

type servantStats struct {
	strength  string
	endurance string
	agility   string
	mana      string
	luck      string
	np        string
}

type servantSkills struct {
	firstSkill    string
	secondSkill   string
	thirdSkill    string
	firstPassive  string
	secondPassive string
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

	wg.Add(7)
	go getServantTags(doc, &wg, s)
	go getServantParams(doc, &wg, s)
	go getServantTraits(doc, &wg, s)
	go getServantCardChoices(doc, &wg, s)
	go getServantCardHits(doc, &wg, s)
	go getServantStatistics(doc, &wg, s)
	go getServantActiveSkills(doc, &wg, s)
	wg.Wait()

	//put all of this in a channel i think
	//split up the tasks concurrently...would be good
	//one function handle alerts
	//one function handle skills
	//one function handle stats etc
}

func getServantTags(doc *goquery.Document, wg *sync.WaitGroup, s *servantCard) {
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

func getServantParams(doc *goquery.Document, wg *sync.WaitGroup, s *servantCard) {
	doc.Find("div[style*='width:300px;margin-top:10px']").First().Find("div[style*='display:inline-block;width:150px;text-align:left;padding: 3px 0 3px 6px;border-left: 2px #4b9acc solid;vertical-align:top']").Each(func(index int, item *goquery.Selection) {
		// use index to determine if jp event only/event character only

		switch index {
		case 0:
			min, max, grail := processParams(item.Text())
			// move all this to the procesParams func
			s.Params.minAtttack = min
			s.Params.maxAtttack = max
			s.Params.grailAttack = grail
		case 1:
			min, max, grail := processParams(item.Text())
			s.Params.minHP = min
			s.Params.maxHP = max
			s.Params.grailHP = grail
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

func getServantTraits(doc *goquery.Document, wg *sync.WaitGroup, s *servantCard) {
	doc.Find("div[style*='width:300px;margin-top:10px;font-size:85%;']").Find("p").Each(func(index int, item *goquery.Selection) {
		s.Traits.traits = item.Text()
	})
	defer wg.Done()
}

func getServantCardChoices(doc *goquery.Document, wg *sync.WaitGroup, s *servantCard) {
	commandQuick := "Command Card Quick.png"
	commandArts := "Command Card Arts.png"
	commandBuster := "Command Card Buster.png"

	quickCount := 0
	artsCount := 0
	busterCount := 0

	doc.Find("div[style*='margin-top: 10px;']").Find("p").Find("img").Each(func(index int, item *goquery.Selection) {
		// temp := item.Attr("width")

		temp, exist := item.Attr("alt")
		if !exist {
			os.Exit(1)
		}

		switch temp {
		case commandQuick:
			quickCount++
			s.CardChoices.quick = quickCount
		case commandArts:
			artsCount++
			s.CardChoices.arts = artsCount
		case commandBuster:
			busterCount++
			s.CardChoices.buster = busterCount
		}
	})

	defer wg.Done()
}

func getServantCardHits(doc *goquery.Document, wg *sync.WaitGroup, s *servantCard) {
	// quickHits := 0
	// artsHits := 0
	// busterHits := 0
	// extraHits := 0

	collectString := doc.Find("div[style*='margin-top: 10px;']").Find("p").Find("span[style*='font-size: 100%']").Find("b").Text()

	tempHold := strings.Split(strings.Replace(removeSymbols(collectString), "NumberofHits", "", -1), "")
	fmt.Println("this is temp hold ", tempHold)

	for i := range tempHold {
		switch i {
		case 0:
			s.CardHits.quickHits = parseStrings(tempHold[i])
		case 1:
			s.CardHits.artsHits = parseStrings(tempHold[i])
		case 2:
			s.CardHits.busterHits = parseStrings(tempHold[i])
		case 3:
			s.CardHits.extraHits = parseStrings(tempHold[i])
		}
	}

	defer wg.Done()
}

func getServantStatistics(doc *goquery.Document, wg *sync.WaitGroup, s *servantCard) {
	doc.Find("div[style*='display:inline-block;vertical-align:top;margin:12px -120px 0 0']").Find("div[style*='display:inline-block;width:150px;text-align:left;padding: 3px 0 3px 6px;border-left: 2px #4b9acc solid;vertical-align:top']").Each(func(index int, item *goquery.Selection) {
		// change values into radar subset here?
		switch index {
		case 0:
			s.Stats.strength = item.Text()
		case 1:
			s.Stats.endurance = item.Text()
		case 2:
			s.Stats.agility = item.Text()
		case 3:
			s.Stats.mana = item.Text()
		case 4:
			s.Stats.luck = item.Text()
		case 5:
			s.Stats.np = item.Text()
		}
		//fmt.Println(item.Text())
	})
	defer wg.Done()
}

// display: block;

func getServantActiveSkills(doc *goquery.Document, wg *sync.WaitGroup, s *servantCard) {
	fmt.Println("servant skills are running")

	doc.Find(".servant-skills").Find("div .tabbertab").Find("div .tabbertab").Each(func(index int, item *goquery.Selection) {

		tempFind, exist := item.Attr("title")
		if exist {
			fmt.Println("this is getServantskills", tempFind)
			fmt.Println("this is the index number", index)
		}

		switch index {
		case 0:
			s.Skills.firstSkill = tempFind
		case 1:
			s.Skills.secondSkill = tempFind
		case 2:
			s.Skills.thirdSkill = tempFind
		}

	})

	defer wg.Done()
}
