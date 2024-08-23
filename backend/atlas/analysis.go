package atlas

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
)

type EconomicData struct {
	Label          string  `json:"label"`
	Value          float64 `json:"value"`
	Risk           float64 `json:"risk"`
	Recommendation string  `json:"recommendation"`
}

type MetricAnalysisText struct {
	Unemployment map[string]string
	Inflation    map[string]string
	InterestRate map[string]string
	GDPGrowth    map[string]string
}

func InitialiseAnalysisTexts() MetricAnalysisText {
	return MetricAnalysisText{
		Unemployment: map[string]string{
			"low":        "Unemployment is significantly below the natural rate, potentially leading to upward wage pressures and inflationary concerns due to a tight labor market.",
			"moderate":   "The unemployment rate is below the natural rate, indicating a strong labor market with minimal slack. Wage growth may accelerate, contributing to inflationary pressures.",
			"medium":     "Unemployment is slightly below equilibrium, suggesting a healthy labor market. However, watch for early signs of labor shortages in key sectors.",
			"balanced":   "Unemployment is near equilibrium, reflecting a balanced labor market. Any significant policy shifts could tip the balance, requiring careful monitoring.",
			"elevated":   "Elevated unemployment levels indicate significant labor market slack, which could necessitate expansionary fiscal or monetary policy interventions.",
			"high":       "Critically high unemployment risk, indicative of severe labor market weakness. Immediate stimulus measures may be required to prevent deflationary spirals.",
			"very_high":  "Unemployment is critically high, indicating a severe economic downturn. Immediate intervention is necessary to avoid a prolonged recession.",
			"extreme":    "Extremely high unemployment, signifying a major economic crisis. Comprehensive and aggressive policy measures are urgently required.",
		},
		Inflation: map[string]string{
			"low":        "Inflation is well within the target range, indicating stable prices. This environment supports sustained economic growth and long-term planning.",
			"moderate":   "Inflation remains under control, though slight upward pressures may be emerging. Policy vigilance is recommended to maintain price stability.",
			"medium":     "Moderate inflationary pressures are beginning to surface, likely due to supply chain constraints or external shocks. A preemptive policy response may be warranted.",
			"balanced":   "Inflation is rising but remains manageable. Continued monitoring and potential fine-tuning of monetary policy could be required to avert further escalation.",
			"elevated":   "High inflation risk, reflecting overheating in the economy. Aggressive monetary tightening may be needed to rein in price growth and anchor expectations.",
			"high":       "Severe inflationary pressures are eroding purchasing power and could destabilize the economy. Coordinated fiscal and monetary actions are urgently required.",
			"very_high":  "Hyperinflation risk is imminent, threatening economic stability. Extreme measures, including potential currency reforms, may be necessary to restore confidence.",
			"extreme":    "Hyperinflation is underway, causing rapid erosion of the currency's value. Immediate and drastic measures are required to stabilize the economy.",
		},
		InterestRate: map[string]string{
			"low":        "Interest rates are at historically low levels, fostering an environment conducive to borrowing and investment. This supports expansionary economic activity.",
			"moderate":   "Interest rates are low, encouraging credit growth and investment. However, potential asset bubbles should be monitored as low rates persist.",
			"medium":     "Interest rates are slightly above the floor, signaling a potential shift towards neutrality. Stakeholders should prepare for possible rate hikes in the near future.",
			"balanced":   "Interest rates are approaching neutrality, suggesting a balanced approach to managing inflation and growth. Market participants should anticipate gradual adjustments.",
			"elevated":   "Elevated interest rates reflect restrictive monetary policy aimed at curbing inflation. The high cost of capital may suppress economic expansion and increase default risks.",
			"high":       "Interest rates are significantly high, suggesting aggressive monetary tightening. The economy could face contractionary pressures as borrowing becomes prohibitively expensive.",
			"very_high":  "Exceptionally high interest rates, likely in response to hyperinflationary threats, could trigger severe economic contraction and destabilize financial markets.",
			"extreme":    "Extremely high interest rates, likely in response to a financial crisis, could lead to a severe economic downturn. Immediate policy intervention is needed.",
		},
		GDPGrowth: map[string]string{
			"low":        "GDP growth is steady and in line with potential output, reflecting a well-balanced economy. Continued prudent policy management is recommended.",
			"moderate":   "GDP growth is moderate, aligning closely with potential output. This suggests stability, though the economy remains vulnerable to external shocks.",
			"medium":     "GDP growth is healthy, slightly above potential output. The economy is performing well, but policymakers should be wary of signs of imbalances.",
			"balanced":   "Strong GDP growth, supported by both domestic and international demand. This growth phase is likely sustainable, though inflationary pressures should be monitored.",
			"elevated":   "GDP growth is exceeding long-term potential, driven by robust demand and favorable external conditions. However, there is a risk of overheating if growth continues unchecked.",
			"high":       "GDP growth is slowing, raising concerns about underlying economic strength. Stimulative measures may be needed to prevent further deceleration.",
			"very_high":  "GDP growth is weak, indicating a decelerating economy. The risk of recession is increasing, requiring proactive counter-cyclical policies.",
			"extreme":    "Critically low GDP growth, signaling a high probability of recession or stagnation. Immediate and significant fiscal and monetary intervention is required to avert a prolonged downturn.",
		},
	}
}


func roundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
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

func (a *Atlas) HandleUnemploymentRate(w http.ResponseWriter, r *http.Request) {
	data, err := a.FetchSeriesData("UNRATE")
	if err != nil || len(data) == 0 {
		http.Error(w, "Failed to fetch unemployment rate data", http.StatusInternalServerError)
		return
	}

	unemploymentRate := roundToTwoDecimals(data[len(data)-1])
	risk := roundToTwoDecimals(min(max((unemploymentRate-4)/6, 0), 1) * 100)
	recommendation := a.GenerateRecommendation("unemployment", risk)

	response := EconomicData{
		Label:          "Unemployment Rate",
		Value:          unemploymentRate,
		Risk:           risk,
		Recommendation: recommendation,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (a *Atlas) HandleInflationRate(w http.ResponseWriter, r *http.Request) {
	data, err := a.FetchSeriesData("CPIAUCSL")
	if err != nil || len(data) == 0 {
		http.Error(w, "Failed to fetch inflation rate data", http.StatusInternalServerError)
		return
	}

	inflationRate := roundToTwoDecimals(((data[len(data)-1] - data[len(data)-13]) / data[len(data)-13]) * 100)
	risk := roundToTwoDecimals(min(max((inflationRate-2)/3, 0), 1) * 100)
	recommendation := a.GenerateRecommendation("inflation", risk)

	response := EconomicData{
		Label:          "Inflation Rate",
		Value:          inflationRate,
		Risk:           risk,
		Recommendation: recommendation,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (a *Atlas) HandleInterestRate(w http.ResponseWriter, r *http.Request) {
	data, err := a.FetchSeriesData("FEDFUNDS")
	if err != nil || len(data) == 0 {
		http.Error(w, "Failed to fetch interest rate data", http.StatusInternalServerError)
		return
	}

	interestRate := roundToTwoDecimals(data[len(data)-1])
	risk := roundToTwoDecimals(min(max(abs(interestRate-2.75)/2.25, 0), 1) * 100)
	recommendation := a.GenerateRecommendation("interest_rate", risk)

	response := EconomicData{
		Label:          "Interest Rate",
		Value:          interestRate,
		Risk:           risk,
		Recommendation: recommendation,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (a *Atlas) HandleGDPGrowth(w http.ResponseWriter, r *http.Request) {
	data, err := a.FetchSeriesData("GDPC1")
	if err != nil || len(data) == 0 {
		http.Error(w, "Failed to fetch GDP growth data", http.StatusInternalServerError)
		return
	}

	gdpGrowth := roundToTwoDecimals(((data[len(data)-1] - data[len(data)-5]) / data[len(data)-5]) * 100)
	risk := roundToTwoDecimals(min(max(abs(gdpGrowth-2.5)/2.5, 0), 1) * 100)
	recommendation := a.GenerateRecommendation("gdp_growth", risk)

	response := EconomicData{
		Label:          "GDP Growth Rate",
		Value:          gdpGrowth,
		Risk:           risk,
		Recommendation: recommendation,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func (a *Atlas) GenerateRecommendation(metric string, risk float64) string {
    analysisTexts := InitialiseAnalysisTexts()

    var analysis map[string]string
    switch metric {
    case "unemployment":
        analysis = analysisTexts.Unemployment
    case "inflation":
        analysis = analysisTexts.Inflation
    case "interest_rate":
        analysis = analysisTexts.InterestRate
    case "gdp_growth":
        analysis = analysisTexts.GDPGrowth
    default:
        return "No specific recommendation available. Further analysis and data collection may be necessary to provide a comprehensive assessment."
    }

    switch {
    case risk < 10:
        return analysis["low"]
    case risk < 20:
        return analysis["moderate"]
    case risk < 30:
        return analysis["medium"]
    case risk < 40:
        return analysis["balanced"]
    case risk < 50:
        return analysis["elevated"]
    case risk < 60:
        return analysis["high"]
    case risk < 70:
        return analysis["very_high"]
    case risk >= 70:
        return analysis["extreme"]
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
