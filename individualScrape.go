package main

import (
	"log"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type servantCard struct {
	EngName  string `json:"engName"`
	JpName   string `json:"jpName"`
	PageLink string `json:"pageLink"`
	Icon     string `json:"icon"`
	Rarity   string `json:"rarity"`
	CardType string `json:"cardType"`
}

func individualScrapeMain(cardType string, pageLink string) {
	if cardType == "servant" {
		scrapeServantPage(pageLink)
	}
}

func scrapeServantPage(pageLink string) servantCard {
	doc, err := goquery.NewDocument(pageLink)
	if err != nil {
		log.Fatal(err)
	}

	// use CSS selector found with the browser inspector
	// for each, use index and item

	var wg sync.WaitGroup
	wg.Add(1)

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
