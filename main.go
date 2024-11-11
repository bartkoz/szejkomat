package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	spreadsheetID = "1uwo-ttmT0DC50z2kENX-2fXK-6qCLF9lKB9sKr1WhaY"
	readRange     = "Sheet1!B2:B"
)

func main() {
	// Set up Google Sheets API client
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile("credentials.json"))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	fmt.Println("Starting to fetch links from Google Sheets...")

	// Retrieve data from the specified range
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	// Check if any data was retrieved
	if len(resp.Values) == 0 {
		fmt.Println("No data found in the specified range.")
		return
	}

	fmt.Println("Successfully fetched list from Google Sheets.")

	// Filter out blank rows and prepare the list of links
	links := []string{}
	for _, row := range resp.Values {
		if len(row) > 0 {
			link := row[0].(string)
			if link != "" {
				links = append(links, link)
			}
		}
	}

	shuffleLinks(links)

	fmt.Printf("Found %d links to process.\n", len(links))

	// Iterate through each link, making GET requests with a delay
	client := &http.Client{}
	for _, link := range links {
		time.Sleep(100 * time.Millisecond) // 0.1 second delay

		// Make GET request to the link
		resp, err := client.Get(link)
		if err != nil {
			fmt.Printf("Failed to GET %s: %v\n", link, err)
			continue
		}

		// Log the response status
		fmt.Printf("GET %s - Status: %s\n", link, resp.Status)
		resp.Body.Close()
	}

	fmt.Println("Finished processing all links.")
}

func shuffleLinks(links []string) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(links), func(i, j int) {
		links[i], links[j] = links[j], links[i]
	})
}
