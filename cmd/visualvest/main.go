package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
	"github.com/subosito/gotenv"
)

type VisualVestTokens struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	IDToken          string `json:"id_token"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

type VisualVestSummary struct {
	Customer struct {
		Gender     string      `json:"gender"`
		Title      interface{} `json:"title"`
		FirstName  string      `json:"firstName"`
		MiddleName string      `json:"middleName"`
		LastName   string      `json:"lastName"`
	} `json:"customer"`
	Hints struct {
		PendingAccountCreation bool          `json:"pendingAccountCreation"`
		MissingDocuments       bool          `json:"missingDocuments"`
		AccountCreationStatus  interface{}   `json:"accountCreationStatus"`
		AccountType            string        `json:"accountType"`
		OfflineMode            bool          `json:"offlineMode"`
		SoldOutDepot           bool          `json:"soldOutDepot"`
		MissingHorizonQuestion bool          `json:"missingHorizonQuestion"`
		ShowBankfusion         bool          `json:"showBankfusion"`
		BankfusionContent      string        `json:"bankfusionContent"`
		LegitimationStatus     interface{}   `json:"legitimationStatus"`
		OptimizationHints      []interface{} `json:"optimizationHints"`
	} `json:"hints"`
	TotalMonthlyRate             float64 `json:"totalMonthlyRate"`
	TotalCurrentAssetValue       float64 `json:"totalCurrentAssetValue"`
	TotalFutureAssetValue        float64 `json:"totalFutureAssetValue"`
	TotalPerformance             float64 `json:"totalPerformance"`
	TotalProfitandlossRealized   float64 `json:"totalProfitandlossRealized"`
	TotalProfitandlossRealizable float64 `json:"totalProfitandlossRealizable"`
	Anlageziele                  []struct {
		Type                          string      `json:"type"`
		ID                            int         `json:"id"`
		Name                          string      `json:"name"`
		ProductType                   interface{} `json:"productType"`
		PortfolioName                 string      `json:"portfolioName"`
		PictureID                     interface{} `json:"pictureId"`
		MonthlyRate                   float64     `json:"monthlyRate"`
		DynamicPercent                int         `json:"dynamicPercent"`
		CurrentAssetValue             float64     `json:"currentAssetValue"`
		FutureAssetValue              float64     `json:"futureAssetValue"`
		Performance                   float64     `json:"performance"`
		ProfitandlossRealized         float64     `json:"profitandlossRealized"`
		ProfitandlossRealizable       float64     `json:"profitandlossRealizable"`
		TargetAmount                  interface{} `json:"targetAmount"`
		TimeFrameYears                interface{} `json:"timeFrameYears"`
		StartDate                     string      `json:"startDate"`
		OptimizationType              interface{} `json:"optimizationType"`
		OptimizationNeedsConfirmation bool        `json:"optimizationNeedsConfirmation"`
		MissingHorizonQuestion        bool        `json:"missingHorizonQuestion"`
		SubDepotNr                    string      `json:"subDepotNr"`
	} `json:"anlageziele"`
}

func init() {
	gotenv.Load()
}

func main() {
	request := gorequest.New()

	// --- 01
	// Login URL
	// ---

	resp, _, errs := request.
		SetDebug(false).
		Get("https://service.visualvest.de/auth/realms/VV/protocol/openid-connect/auth").
		Param("client_id", "dashboard-login-switch-frontend").
		Param("redirect_uri", "https://geldanlage.visualvest.de/login-switch/").
		Param("response_type", "code").
		End()
	if errs != nil {
		fmt.Println(errs)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	loginForm := doc.Find("#kc-form-login")
	loginURL, exists := loginForm.Attr("action")
	if exists != true {
		panic("no login URL found")
	}

	// --- 02
	// Login
	// ---

	resp, _, errs = request.
		SetDebug(false).
		Post(loginURL).
		Send("username=perelin&password=Z2GHOw%25y%2655e&login=").
		RedirectPolicy(func(req gorequest.Request, via []gorequest.Request) error {
			return http.ErrUseLastResponse
		}).
		End()
	if errs != nil {
		fmt.Println(errs)
	}

	// --- 03
	// Auth against dashboard-frontend-fpv
	// ---

	resp, _, errs = request.
		SetDebug(false).
		Get("https://service.visualvest.de/auth/realms/VV/protocol/openid-connect/auth").
		Param("client_id", "dashboard-frontend-fpv").
		Param("redirect_uri", "https://geldanlage.visualvest.de/dashboard").
		Param("response_type", "code").
		End()
	if errs != nil {
		fmt.Println(errs)
	}

	parsedURL, err := url.Parse(resp.Header["Location"][0])
	if err != nil {
		panic(err)
	}
	m, _ := url.ParseQuery(parsedURL.RawQuery)
	//fmt.Println(m["code"][0])

	// --- 04
	// Token
	// ---

	var tokens VisualVestTokens

	resp, _, errs = request.
		SetDebug(false).
		Post("https://service.visualvest.de/auth/realms/VV/protocol/openid-connect/token").
		Send("code=" + m["code"][0] + "&grant_type=authorization_code&client_id=dashboard-frontend-fpv&redirect_uri=https%3A%2F%2Fgeldanlage.visualvest.de%2Fdashboard").
		EndStruct(&tokens)
	if errs != nil {
		fmt.Println(errs)
	}
	//fmt.Println(tokens)

	// --- 05
	// Summary
	// ---

	var summary VisualVestSummary
	resp, _, errs = request.
		SetDebug(false).
		Get("https://service-geldanlage.visualvest.de/investment-summary-functional-service/investment-summary/depot/1").
		Set("Authorization", "Bearer "+tokens.AccessToken).
		EndStruct(&summary)
	if errs != nil {
		fmt.Println(errs)
	}
	fmt.Println(summary.TotalCurrentAssetValue)
	spew.Dump(summary)

	// --- POST
	// Write financial overview to google sheets
	// ---

	srv := getSheetsService()
	now := time.Now()
	values := []interface{}{now.Format("02.01.2006 15:04:05"),
		summary.TotalCurrentAssetValue,
		summary.TotalPerformance,
	}
	appendToOnvistSheet(srv, values)

}
