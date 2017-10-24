package main

import (
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type servant struct {
	EngName  string `json:"engName"`
	JpName   string `json:"jpName"`
	PageLink string `json:"pageLink"`
	Icon     string `json:"icon"`
	Rarity   string `json:"rarity"`
	CardType string `json:"cardType"`
}

func servantsMain() []servant {
	// measure entire function
	start := time.Now()

	r := new(big.Int)
	fmt.Println(r.Binomial(1000, 10))

	var servantCollection []servant

	servantLength := findServantLength()
	servantChannel := make(chan servant, servantLength)
	tempServant := scrapeServantNames(servantChannel, servantLength)
	close(servantChannel)

	// use a merge go routine to combine, look at factorial redux
	// turn channel into a slice
	for n := range tempServant {
		if len(n.JpName) > 0 {
			// fmt.Println(n)
			servantCollection = append(servantCollection, n)
		}
	}

	elapsed := time.Since(start)
	log.Printf("Servant search took %s", elapsed)

	return servantCollection
}

func findServantLength() int {
	doc, err := goquery.NewDocument("https://grandorder.wiki/Servant_List")
	if err != nil {
		log.Fatal(err)
	}

	var lengthDoc = doc.Find("tr").Length()

	return lengthDoc
}

func scrapeServantNames(s chan servant, length int) chan servant {
	doc, err := goquery.NewDocument("https://grandorder.wiki/Servant_List")
	if err != nil {
		log.Fatal(err)
	}

	// use CSS selector found with the browser inspector
	// for each, use index and item

	var wg sync.WaitGroup
	wg.Add(length)

	doc.Find("tr").Each(func(index int, item *goquery.Selection) {
		// title := strings.TrimSpace(item.Text())
		// linkTag := item.Find("a")
		// link, _ := linkTag.Attr("href")
		// fmt.Printf("Post #%d: %s - %s\n", index, title, link)

		engName := item.Find("td").First().Next().Next().Find("a").Text()
		jpName := item.Find("td").First().Next().Next().Find("small").Text()

		servantLink, _ := item.Find("td").First().Next().Next().Find("a").Attr("href")
		servantLink = "https://grandorder.wiki" + servantLink

		imageLink, _ := item.Find("td").First().Next().Find("img").Attr("src")
		imageLink = "https://grandorder.wiki" + imageLink

		cardType := "servant"

		rarity := item.Find("td").First().Next().Next().Next().Next().Text()

		go func(s chan servant) {
			defer wg.Done()
			tempServant := servant{
				EngName:  engName,
				JpName:   jpName,
				PageLink: servantLink,
				Icon:     imageLink,
				Rarity:   rarity,
				CardType: cardType,
			}
			s <- tempServant
		}(s)
	})
	wg.Wait()
	return s
}
