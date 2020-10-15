package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector()
	fmt.Print("Visiting http://go-colly.org/... ")
	c.Visit("http://go-colly.org/")
	fmt.Println("Done")
}
