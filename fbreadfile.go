package main

import (
	"io/ioutil"
	"log"
	"strings"

	"fmt"

	"golang.org/x/net/html"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getVal(t html.Token, key string) (ok bool, val string) {
	// Iterate over all of the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == key {
			val = a.Val
			ok = true
		}
	}
	return
}

func extractText(s string) (texts []string) {
	var wasNl = false
	z := html.NewTokenizer(strings.NewReader(s))
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			//log.Println("err tok")
			return
		case tt == html.TextToken:
			t := z.Token()
			text := strings.TrimSpace(t.Data)
			if text != "" {
				texts = append(texts, strings.TrimSpace(t.Data))
				wasNl = false
				log.Println("text:", "'"+strings.TrimSpace(t.Data)+"'")
			}

		case tt == html.StartTagToken || tt == html.SelfClosingTagToken:
			t := z.Token()
			switch t.Data {
			case "br":
				if !wasNl {
					texts = append(texts, "\n")
				}
				wasNl = true
			}
		}
	}
}

func extractLinks(s string) (links []string) {

	z := html.NewTokenizer(strings.NewReader(s))
	var username bool
	var dolj bool
	var printlink bool
	var i = 0
	for {
		tt := z.Next()
		//log.Println("tt", tt)
		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			//log.Println("err tok")
			return
		case tt == html.TextToken:
			t := z.Token()
			text := strings.TrimSpace(t.Data)
			if username {
				i++
				fmt.Printf("%d\t%s", i, text)
				username = false
			}
			if dolj {
				fmt.Printf("\t%s", text)
				dolj = false
				printlink = true
			}
			//log.Println(text)
		case tt == html.StartTagToken || tt == html.SelfClosingTagToken:
			t := z.Token()
			//log.Println("t", t)
			switch t.Data {
			case "div":
				ok, name := getVal(t, "class")
				if !ok {
					continue
				}
				switch name {
				case "_5d-5":
					username = true
				case "_pac":
					dolj = true
				}

			case "img":
				//log.Println("img")
				ok, href := getVal(t, "src")
				if !ok {
					continue
				}
				//checkLink(href)

				_ = href
			case "a":
				ok, href := getVal(t, "href")
				if !ok {
					continue
				}
				if printlink {
					fmt.Printf("\t%s\n", href)
					printlink = false
				}

			default:
				//log.Println("t.data", t.Data)
				continue
			}
		}
	}
}

func extractIds(s string) (links []string) {

	z := html.NewTokenizer(strings.NewReader(s))
	var i = 0
	var lastid = ""
	for {
		tt := z.Next()
		//log.Println("tt", tt)
		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			//log.Println("err tok")
			return
		case tt == html.StartTagToken || tt == html.SelfClosingTagToken:
			t := z.Token()
			//log.Println("t", t)
			switch t.Data {

			case "a":
				ok, href := getVal(t, "href")
				if !ok {
					continue
				}
				ok, keys := getVal(t, "data-testid")
				if !ok {
					continue
				}
				//log.Println(keys)
				if strings.Contains(keys, "EntRegularPersonalUser") {
					if lastid == href {
						continue
					}
					lastid = href
					i++
					fmt.Printf("%d\t%s\n", i, href)
				}

			default:
				//log.Println("t.data", t.Data)
				continue
			}
		}
	}
}

func main() {

	//https://www.facebook.com/search/str/mail%20ru%20group/pages-named/employees/present/intersect
	//shalaev.s
	//1385010678
	//emaleev
	//533349026
	dat, err := ioutil.ReadFile("4.html")
	check(err)
	s := string(dat)
	//extractLinks(s)
	extractIds(s)

}
