package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
	"github.com/guimassoqueto/go-fake-headers"
)

func main() {
	urls := []string{
		"https://amazon.com.br/dp/B08R6KNV45",
		"https://amazon.com.br/dp/B0B8YT1X8H",
		"https://amazon.com.br/dp/B07ZD8RR5C",
		"https://amazon.com.br/dp/B086F8QWW5",
		"https://amazon.com.br/dp/B000PQECW6",
		"https://amazon.com.br/dp/B07HDSL78H",
		"https://amazon.com.br/dp/B005IZZDAY",
		"https://amazon.com.br/dp/B07M73V8JD",
		"https://amazon.com.br/dp/B07XLMJBWQ",
		"https://amazon.com.br/dp/B00H301K6I",
		"https://amazon.com.br/dp/B07YN9M2MH",
		"https://amazon.com.br/dp/B0BVZY7BM6",
		"https://amazon.com.br/dp/B07X67Z1RP",
		"https://amazon.com.br/dp/B07YN9ZY5D",
		"https://amazon.com.br/dp/B07XMMMQ38",
		"https://amazon.com.br/dp/B07MC56KN3",
		"https://amazon.com.br/dp/B0B6QB1R86",
		"https://amazon.com.br/dp/B00FOYOE1S",
		"https://amazon.com.br/dp/B09M2B14V7",
		"https://amazon.com.br/dp/B0C6WDYV6X",
		"https://amazon.com.br/dp/B07TVFQVJN",
		"https://amazon.com.br/dp/B0823BN8FR",
		"https://amazon.com.br/dp/B075W7XHMY",
		"https://amazon.com.br/dp/B07DNJ2QSP",
		"https://amazon.com.br/dp/B00MU6E8BY",
		"https://amazon.com.br/dp/B07YN9SQR9",
		"https://amazon.com.br/dp/B097BYXGXN",
		"https://amazon.com.br/dp/B08239Y8N1",
	}

	// Create a WaitGroup to wait for all Goroutines to finish
	var wg sync.WaitGroup

	// Create a channel for the worker pool with a capacity of 16
	concurrentLinks := make(chan string, 32)

	// Create a new Colly collector
	c := colly.NewCollector(
		colly.AllowedDomains(),
		colly.IgnoreRobotsTxt(), // Set IgnoreRobotsTxt to true to avoid tracking visited URLs
	)

	c.OnRequest(func(r *colly.Request) {
		fakeHeader := randomHeader.Build()
		for k, v := range fakeHeader {
			r.Headers.Set(k, v)
		}
	})
	
	var title string
	c.OnHTML("#title", func(e *colly.HTMLElement) {
		// Extract the title of the webpage
		title = strings.Trim(e.Text, " ")
		fmt.Println(title)
		wg.Done() // Mark the Goroutine as done when processing is complete
	})

	// Handle errors
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		wg.Done() // Mark the Goroutine as done even if there was an error
	})

	// Start Goroutines for each URL
	for _, url := range urls {
		wg.Add(1) // Increment the WaitGroup counter for each URL

		// Send the URL to the worker pool for processing
		concurrentLinks <- url

		go func(u string) {
			// Defer the removal of the URL from the worker pool
			defer func() { <-concurrentLinks }()

			// Make the request using Colly
			err := c.Visit(u)
			if err != nil {
				fmt.Println("Error visiting URL:", u, "\nError:", err)
				return
			}
		}(url)
	}

	// Wait for all Goroutines to finish
	wg.Wait()
	close(concurrentLinks)
}