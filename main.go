package main

import (
	"fmt"
	"os"

	"KumaCapital/atlas"
	"github.com/gin-gonic/gin"
)

func main() {
	apiKey := os.Getenv("FRED_API_KEY")
	if apiKey == "" {
		fmt.Println("FRED_API_KEY environment variable is not set")
		return
	}

	atlasAnalyser := atlas.InitialiseAtlas(apiKey)

	router := gin.Default()

	router.GET("/atlas/unemployment", atlasAnalyser.AnalyseUnemploymentRate)
	router.GET("/atlas/inflation", atlasAnalyser.AnalyseInflationRate)
	router.GET("/atlas/interest_rate", atlasAnalyser.AnalyseInterestRate)
	router.GET("/atlas/gdp_growth", atlasAnalyser.AnalyseGDPGrowth)

	router.Run(":8080")
}
