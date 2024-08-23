package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"KumaCapital/atlas"
)

func main() {
	apiKey := os.Getenv("FRED_API_KEY")
	if apiKey == "" {
		log.Fatal("FRED_API_KEY environment variable is not set")
	}

	atlasAnalyser := atlas.InitialiseAtlas(apiKey)

	http.HandleFunc("/atlas/unemployment", atlasAnalyser.HandleUnemploymentRate)
	http.HandleFunc("/atlas/inflation", atlasAnalyser.HandleInflationRate)
	http.HandleFunc("/atlas/interest_rate", atlasAnalyser.HandleInterestRate)
	http.HandleFunc("/atlas/gdp_growth", atlasAnalyser.HandleGDPGrowth)

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
