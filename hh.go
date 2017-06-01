package main

import (
	"fmt"
	"log"
	"os"

	"strings"

	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func cleanText(s string) (res string) {
	res = strings.Replace(s, "\n", "", -1)
	res = strings.TrimSpace(res)
	return
}

func Scrape(name string) {
	file, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		log.Fatal(err)
	}

	table := doc.Find("table .output")
	table.Each(func(_ int, tables *goquery.Selection) {
		trs := tables.Find("tr")
		trs.Each(func(i int, tr *goquery.Selection) {
			hash, _ := tr.Attr("data-hh-resume-hash")
			title := tr.Find("a").Text()
			opit := tr.Find("div .output__experience-sum").Text()
			last := tr.Find("div .output__indent strong").Text()
			age := cleanText(tr.Find(".output__age").Text())
			comp := cleanText(tr.Find(".output__compensation").Text())
			fmt.Printf("%d	%s	%s	%s	%s	%s	%s\n", i, "https://hh.ru/resume/"+hash, age, comp, title, opit, last)
		})
	})
}

func main() {
	fmt.Printf("i	resume	age	comp	title	opit	last\n")
	for i := 1; i < 17; i++ {
		Scrape("HeadHunter" + strconv.Itoa(i) + ".html")
	}

}
