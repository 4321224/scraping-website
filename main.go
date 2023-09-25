package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
    "gopkg.in/mgo.v2"
)

type Product struct {
	Name  string `bson:"name"`
	Price string `bson:"price"`
}

func ScrapeZara() []Product {
	var products []Product

	res, err := http.Get("https://www.zara.com/id/")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".product").Each(func(i int, s *goquery.Selection) {
        name := s.Find(".name").Text()
        price := s.Find(".price").Text()

        product := Product{
            Name:  name,
            Price: price,
        }

        products = append(products, product)
    })

	return products
}

func main() {
	session, err := mgo.Dial("mongodb://localhost:27017/")
	if err != nil {
	    panic(err)
    }
	defer session.Close()

	collection := session.DB("test").C("zara")

    products := ScrapeZara()
	for _, product := range products {
	    err = collection.Insert(&product)

	    if err != nil {
	        log.Fatal(err)
	    }
    }

	fmt.Println("Data scraped and stored successfully!")
}
