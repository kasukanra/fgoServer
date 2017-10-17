package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type craftEssence struct {
	EngName  string `json:"engName"`
	PageLink string `json:"pageLink"`
	Icon     string `json:"icon"`
	Rarity   string `json:"rarity"`
}

func ceMain() []craftEssence {
	var ceCollection []craftEssence

	ceLength := findCELength()
	ceChannel := make(chan craftEssence, ceLength)
	tempCE := scrapeCENames(ceChannel, ceLength)
	close(ceChannel)

	// use a merge go routine to combine, look at factorial redux
	// turn channel into a slice
	for n := range tempCE {
		if len(n.EngName) > 0 {
			fmt.Println(n)
			ceCollection = append(ceCollection, n)
		}
	}

	return ceCollection
}

func findCELength() int {
	doc, err := goquery.NewDocument("https://grandorder.wiki/Craft_Essences")
	if err != nil {
		log.Fatal(err)
	}

	var lengthDoc = doc.Find("tr").Length()

	return lengthDoc
}

func scrapeCENames(c chan craftEssence, length int) chan craftEssence {
	doc, err := goquery.NewDocument("https://grandorder.wiki/Craft_Essences")
	if err != nil {
		log.Fatal(err)
	}

	// use CSS selector found with the browser inspector
	// for each, use index and item

	var wg sync.WaitGroup
	wg.Add(length)

	doc.Find("tr").Each(func(index int, item *goquery.Selection) {
		engName := item.Find("td").First().Next().Next().Find("a").Text()

		fmt.Println("engName", engName)

		ceLink, _ := item.Find("td").First().Next().Next().Find("a").Attr("href")
		ceLink = "https://grandorder.wiki" + ceLink

		imageLink, _ := item.Find("td").First().Next().Find("img").Attr("src")
		imageLink = "https://grandorder.wiki" + imageLink

		rarity := item.Find("td").First().Next().Next().Next().Text()

		// arr = append(arr, engName)

		// fmt.Println(engName)
		// fmt.Println(jpName)

		go func(c chan craftEssence) {
			defer wg.Done()
			tempCE := craftEssence{
				EngName:  engName,
				PageLink: ceLink,
				Icon:     imageLink,
				Rarity:   rarity,
			}
			c <- tempCE
		}(c)

	})
	wg.Wait()
	return c
}
