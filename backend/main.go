package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/rs/cors"
    "KumaCapital/atlas"
)

func main() {
    apiKey := os.Getenv("FRED_API_KEY")
    if apiKey == "" {
        log.Fatal("FRED_API_KEY environment variable is not set")
    }

    atlasAnalyser := atlas.InitialiseAtlas(apiKey)

    mux := http.NewServeMux()
    mux.HandleFunc("/atlas/unemployment", atlasAnalyser.HandleUnemploymentRate)
    mux.HandleFunc("/atlas/inflation", atlasAnalyser.HandleInflationRate)
    mux.HandleFunc("/atlas/interest_rate", atlasAnalyser.HandleInterestRate)
    mux.HandleFunc("/atlas/gdp_growth", atlasAnalyser.HandleGDPGrowth)

    handler := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000"}, 
        AllowCredentials: true,
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type"},
    }).Handler(mux)

    fmt.Println("Server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", handler))
}
