package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {

	url := "https://shampion.mk"

	resp, err := http.Get(url) //zema response ili greshka od sajto sho mu go davame
	if err != nil {
		log.Fatal("URL e nedostapno: ", err)
	}
	defer resp.Body.Close() //ova go zatvora responese body od ko ce procita za da ne bara dodatno

	//sledno proveruva dali responso ni e 200 ako ne e frla exception
	if resp.StatusCode != 200 {
		log.Fatal("Greska pri vcituvanje na kodo", err)
	}

	//koristam goquery za da go zemi HTML-o, ako ne mozhi frla exception
	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		log.Fatal("Ne mozi da se vcita HTML:", err)
	}

	//funkcija za da isfiltrira se sho ne e tekst

	filterNonText := func(i int, s *goquery.Selection) bool {
		tagName := strings.ToLower(s.Get(0).Data)
		return tagName == "p" || tagName == "div" || tagName == "span" || tagName == "a"
		//vrajcha p, div, span i a sodrzhini
	}
	//Sledno bara vo body preku koristenje na filter funkcijata so koristenje kako vlezen parametar funkcijata za filtriranje na non-text
	doc.Find("body *").FilterFunction(filterNonText).Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			fmt.Println(text)
		}
	})

	//	ovde dolu gi trga nepotrebnite kodovi, vo slucajov css javascript i jquery za da se dobija chist tekst

	doc.Find("style").Remove()
	doc.Find("script").Remove()
	doc.Find(".jquery-script").Remove()

	lastWasSpace := true // che proveruva dali poslednio element imal mesto
	var textContent strings.Builder
	doc.Find("body *").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {

			if !lastWasSpace && textContent.Len() > 0 {
				textContent.WriteString(" ") //da dodaj mesto izmedzu ako nema voopsto
			}
			textContent.WriteString(text)
			lastWasSpace = false
		} else {
			lastWasSpace = true
		}
	})

	fmt.Println("TEXT THAT I NEED:")
	fmt.Println(textContent.String())

	// da se zacuvuva seto ova vo text
	fileName := "extracted_content.txt"
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Greshka pri kreiranje na fajlo: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(textContent.String())
	if err != nil {
		log.Fatalf("Greshka pri writing na fajlo %v", err)
	}

	//ovde samo printa dali se zacuvuva
	fmt.Printf("fajlot e zacuvan %s\n", fileName)

	absPath, err := filepath.Abs(fileName)
	if err != nil {
		log.Fatalf("Ne mozi da se dobija patekata  %v", err)
	}
	fmt.Printf("Teksto e zacuvan vo:\n%s\n", absPath)

}
