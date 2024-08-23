package atlas

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EconomicData struct {
	Label        string  `json:"label"`
	Value        float64 `json:"value"`
	Risk         float64 `json:"risk"`
	Recommendation string `json:"recommendation"`
}

func (a *Atlas) FetchSeriesData(seriesCode string) ([]float64, error) {
	resp, err := a.Client.R().
		SetQueryParams(map[string]string{
			"series_id": seriesCode,
			"api_key":   a.ApiKey,
			"file_type": "json",
		}).
		Get("https://api.stlouisfed.org/fred/series/observations")

	if err != nil {
		return nil, err
	}

	var result struct {
		Observations []struct {
			Value string `json:"value"`
		} `json:"observations"`
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}

	var data []float64
	for _, observation := range result.Observations {
		val, err := strconv.ParseFloat(observation.Value, 64)
		if err != nil {
			return nil, err
		}
		data = append(data, val)
	}

	return data, nil
}

func (a *Atlas) AnalyseUnemploymentRate(c *gin.Context) {
	data, err := a.FetchSeriesData("UNRATE")
	if err != nil || len(data) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch unemployment rate data"})
		return
	}

	unemploymentRate := data[len(data)-1]
	risk := min(max((unemploymentRate-4)/6, 0), 1) * 100
	recommendation := a.GenerateRecommendation("unemployment", risk)

	response := EconomicData{
		Label:        "Unemployment Rate",
		Value:        unemploymentRate,
		Risk:         risk,
		Recommendation: recommendation,
	}

	c.JSON(http.StatusOK, response)
}

func (a *Atlas) AnalyseInflationRate(c *gin.Context) {
	data, err := a.FetchSeriesData("CPIAUCSL")
	if err != nil || len(data) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch inflation rate data"})
		return
	}

	inflationRate := ((data[len(data)-1] - data[len(data)-13]) / data[len(data)-13]) * 100
	risk := min(max((inflationRate-2)/3, 0), 1) * 100
	recommendation := a.GenerateRecommendation("inflation", risk)

	response := EconomicData{
		Label:        "Inflation Rate",
		Value:        inflationRate,
		Risk:         risk,
		Recommendation: recommendation,
	}

	c.JSON(http.StatusOK, response)
}

func (a *Atlas) AnalyseInterestRate(c *gin.Context) {
	data, err := a.FetchSeriesData("FEDFUNDS")
	if err != nil || len(data) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch interest rate data"})
		return
	}

	interestRate := data[len(data)-1]
	risk := min(max(abs(interestRate-2.75)/2.25, 0), 1) * 100
	recommendation := a.GenerateRecommendation("interest_rate", risk)

	response := EconomicData{
		Label:        "Interest Rate",
		Value:        interestRate,
		Risk:         risk,
		Recommendation: recommendation,
	}

	c.JSON(http.StatusOK, response)
}

func (a *Atlas) AnalyseGDPGrowth(c *gin.Context) {
	data, err := a.FetchSeriesData("GDPC1")
	if err != nil || len(data) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch GDP growth data"})
		return
	}

	gdpGrowth := ((data[len(data)-1] - data[len(data)-5]) / data[len(data)-5]) * 100
	risk := min(max(abs(gdpGrowth-2.5)/2.5, 0), 1) * 100
	recommendation := a.GenerateRecommendation("gdp_growth", risk)

	response := EconomicData{
		Label:        "GDP Growth Rate",
		Value:        gdpGrowth,
		Risk:         risk,
		Recommendation: recommendation,
	}

	c.JSON(http.StatusOK, response)
}

func (a *Atlas) GenerateRecommendation(metric string, risk float64) string {
	switch metric {
	case "unemployment":
		if risk < 10 {
			return "Unemployment is significantly below the natural rate, potentially leading to upward wage pressures and inflationary concerns due to a tight labor market."
		} else if risk < 20 {
			return "The unemployment rate is below the natural rate, indicating a strong labor market with minimal slack. Wage growth may accelerate, contributing to inflationary pressures."
		} else if risk < 30 {
			return "Unemployment is slightly below equilibrium, suggesting a healthy labor market. However, watch for early signs of labor shortages in key sectors."
		} else if risk < 40 {
			return "Unemployment is near equilibrium, reflecting a balanced labor market. Any significant policy shifts could tip the balance, requiring careful monitoring."
		} else if risk < 50 {
			return "Unemployment is slightly above equilibrium, indicating emerging slack in the labor market. This could signal the early stages of an economic slowdown."
		} else if risk < 60 {
			return "Rising unemployment suggests increasing labor market slack, potentially leading to reduced consumer spending and a drag on economic growth."
		} else if risk < 70 {
			return "Elevated unemployment levels indicate significant labor market slack, which could necessitate expansionary fiscal or monetary policy interventions."
		} else {
			return "Critically high unemployment risk, indicative of severe labor market weakness. Immediate stimulus measures may be required to prevent deflationary spirals."
		}
	case "inflation":
		if risk < 10 {
			return "Inflation is well within the target range, indicating stable prices. This environment supports sustained economic growth and long-term planning."
		} else if risk < 20 {
			return "Inflation remains under control, though slight upward pressures may be emerging. Policy vigilance is recommended to maintain price stability."
		} else if risk < 30 {
			return "Moderate inflationary pressures are beginning to surface, likely due to supply chain constraints or external shocks. A preemptive policy response may be warranted."
		} else if risk < 40 {
			return "Inflation is rising but remains manageable. Continued monitoring and potential fine-tuning of monetary policy could be required to avert further escalation."
		} else if risk < 50 {
			return "Inflationary risks are becoming more pronounced, driven by persistent demand-side pressures or supply shortages. Consideration of policy tightening is advisable."
		} else if risk < 60 {
			return "High inflation risk, reflecting overheating in the economy. Aggressive monetary tightening may be needed to rein in price growth and anchor expectations."
		} else if risk < 70 {
			return "Severe inflationary pressures are eroding purchasing power and could destabilise the economy. Coordinated fiscal and monetary actions are urgently required."
		} else {
			return "Hyperinflation risk is imminent, threatening economic stability. Extreme measures, including potential currency reforms, may be necessary to restore confidence."
		}
	case "interest_rate":
		if risk < 10 {
			return "Interest rates are at historically low levels, fostering an environment conducive to borrowing and investment. This supports expansionary economic activity."
		} else if risk < 20 {
			return "Interest rates are low, encouraging credit growth and investment. However, potential asset bubbles should be monitored as low rates persist."
		} else if risk < 30 {
			return "Interest rates are slightly above the floor, signaling a potential shift towards neutrality. Stakeholders should prepare for possible rate hikes in the near future."
		} else if risk < 40 {
			return "Interest rates are approaching neutrality, suggesting a balanced approach to managing inflation and growth. Market participants should anticipate gradual adjustments."
		} else if risk < 50 {
			return "Interest rates are moderately high, indicating a cautious approach to inflation control. Borrowing costs are rising, potentially dampening investment and consumption."
		} else if risk < 60 {
			return "Elevated interest rates reflect restrictive monetary policy aimed at curbing inflation. The high cost of capital may suppress economic expansion and increase default risks."
		} else if risk < 70 {
			return "Interest rates are significantly high, suggesting aggressive monetary tightening. The economy could face contractionary pressures as borrowing becomes prohibitively expensive."
		} else {
			return "Exceptionally high interest rates, likely in response to hyperinflationary threats, could trigger severe economic contraction and destabilise financial markets."
		}
	case "gdp_growth":
		if risk < 10 {
			return "GDP growth is exceeding long-term potential, driven by robust demand and favorable external conditions. However, there is a risk of overheating if growth continues unchecked."
		} else if risk < 20 {
			return "Strong GDP growth, supported by both domestic and international demand. This growth phase is likely sustainable, though inflationary pressures should be monitored."
		} else if risk < 30 {
			return "GDP growth is healthy, slightly above potential output. The economy is performing well, but policymakers should be wary of signs of imbalances."
		} else if risk < 40 {
			return "GDP growth is steady and in line with potential output, reflecting a well-balanced economy. Continued prudent policy management is recommended."
		} else if risk < 50 {
			return "GDP growth is moderate, aligning closely with potential output. This suggests stability, though the economy remains vulnerable to external shocks."
		} else if risk < 60 {
			return "GDP growth is slowing, raising concerns about underlying economic strength. Stimulative measures may be needed to prevent further deceleration."
		} else if risk < 70 {
			return "GDP growth is weak, indicating a decelerating economy. The risk of recession is increasing, requiring proactive counter-cyclical policies."
		} else {
			return "Critically low GDP growth, signaling a high probability of recession or stagnation. Immediate and significant fiscal and monetary intervention is required to avert a prolonged downturn."
		}
	default:
		return "No specific recommendation available. Further analysis and data collection may be necessary to provide a comprehensive assessment."
	}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func abs(a float64) float64 {
	if a < 0 {
		return -a
	}
	return a
}
