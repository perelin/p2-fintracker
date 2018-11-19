package main

import (
	"fmt"
	"os"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/subosito/gotenv"
)

type OnvistaFinancialOverview struct {
	S0 struct {
		Result struct {
			FinancialOverview struct {
				Depot struct {
					BuyValue                float64 `json:"buyValue"`
					ActualValue             float64 `json:"actualValue"`
					TotalPerformance        float64 `json:"totalPerformance"`
					PerformancePercentage   float64 `json:"performancePercentage"`
					TotalDailyPerformance   float64 `json:"totalDailyPerformance"`
					TotalDailyPerformancePx float64 `json:"totalDailyPerformancePx"`
				} `json:"depot"`
			} `json:"financialOverview"`
			Meta struct {
				RequestExecutionTime float64 `json:"requestExecutionTime"`
			} `json:"_meta"`
		} `json:"result"`
	} `json:"s0"`
}

func init() {
	gotenv.Load()
}

func main() {

	// --- AUTH
	// Authenticate against onvista
	// ---

	request := gorequest.New()
	_, _, errs := request.SetDebug(false).
		Post("https://webtrading.onvista-bank.de/services/api/?s0=Session_Auth.login").
		Type("multipart").
		Send("hash[key]=JKEIMG1J1NIJ5619").
		Send("action[s0][domain]=Session_Auth").
		Send("action[s0][service]=login").
		Send("action[s0][params][login]=" + os.Getenv("ONVISTA_USER")).
		Send("action[s0][params][password]=" + os.Getenv("ONVISTA_PWD")).
		End()
	if errs != nil {
		fmt.Println(errs)
	}

	// --- GET
	// Get onvista financial overview
	// ---

	var financialOverview OnvistaFinancialOverview
	_, _, errs = request.SetDebug(false).
		Post("https://webtrading.onvista-bank.de/services/api/?s0=Session_Auth.login").
		Type("multipart").
		Send("hash[key]=JKEIMG1J1NIJ5619").
		Send("action[s0][domain]=Bank_Overview").
		Send("action[s0][service]=getFinancialOverview").
		Send("action[s0][params][accountKey]=a2847912e0e52940bfd8bcb19e964cfa").
		EndStruct(&financialOverview)
	if errs != nil {
		fmt.Println(errs)
	}

	// --- POST
	// Write financial overview to google sheets
	// ---

	if financialOverview.S0.Result.FinancialOverview.Depot.BuyValue == 0 {
		// resilience against nonsense API return values
		return
	}

	srv := getSheetsService()
	now := time.Now()
	values := []interface{}{now.Format("02.01.2006 15:04:05"),
		financialOverview.S0.Result.FinancialOverview.Depot.BuyValue,
		financialOverview.S0.Result.FinancialOverview.Depot.ActualValue,
		financialOverview.S0.Result.FinancialOverview.Depot.TotalPerformance,
		financialOverview.S0.Result.FinancialOverview.Depot.PerformancePercentage / 100,
		financialOverview.S0.Result.FinancialOverview.Depot.TotalDailyPerformance,
		financialOverview.S0.Result.FinancialOverview.Depot.TotalDailyPerformancePx / 100,
	}
	appendToOnvistSheet(srv, values)

}
