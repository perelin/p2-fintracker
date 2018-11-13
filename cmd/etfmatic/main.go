package main

import (
	"fmt"
	"time"

	"github.com/parnurzeal/gorequest"
)

type emSignin struct {
	Success    bool `json:"success"`
	ClientMode struct {
		ID                 string `json:"id"`
		Name               string `json:"name"`
		Description        string `json:"description"`
		ClientLoginAllowed bool   `json:"client_login_allowed"`
	} `json:"client_mode"`
	Token         string `json:"token"`
	ClientSystems struct {
		EuLive     string `json:"eu_live"`
		Simulation string `json:"simulation"`
	} `json:"client_systems"`
}

func main() {
	// --- AUTH
	// Authenticate against onvista
	// ---

	request := gorequest.New()
	var emSignin emSignin
	_, _, errs := request.SetDebug(false).
		Post("https://gateway.etfmatic.com/user/1.0/sign-in").
		Type("multipart").
		Send("api_data={\"user_email\":\"patino@p2lab.de\",\"user_password\":\"go17:ab03\"}").EndStruct(&emSignin)
	if errs != nil {
		fmt.Println(errs)
	}

	//spew.Dump(emSignin.Token)
	//spew.Dump(emSignin.ClientSystems.EuLive)

	// var emAll emAll
	// _, _, errs = request.
	// 	SetDebug(false).
	// 	Get("https://gateway.etfmatic.com/user/1.0/retrieve").
	// 	Param("session_token", emSignin.Token).
	// 	Param("system_identifier", emSignin.ClientSystems.EuLive).
	// 	EndStruct(&emAll)

	// if errs != nil {
	// 	fmt.Println(errs)
	// }
	// spew.Dump(emAll.UserAccountsSummary.TotalContributions)

	var emGoals emGoals
	_, _, errs = request.
		SetDebug(false).
		Get("https://gateway.etfmatic.com/user-goal/1.0/retrieve-goals").
		Param("session_token", emSignin.Token).
		Param("system_identifier", emSignin.ClientSystems.EuLive).
		EndStruct(&emGoals)
	if errs != nil {
		fmt.Println(errs)
	}

	// spew.Dump(emGoals.UsrGoals[0].UsrinvCurrentValuation)
	// spew.Dump(emGoals.UsrGoals[0].UsrinvCurrentValuationAssets)
	// spew.Dump(emGoals.UsrGoals[0].UsrinvTotalGains)
	// spew.Dump(emGoals.UsrGoals[0].UsrinvTotalDividendsReceived)
	// spew.Dump(emGoals.UsrGoals[0].UsrinvTotalContributions)
	// spew.Dump(emGoals.UsrGoals[0].UsrinvUnallocatedCreditAmount)

	srv := getSheetsService()
	now := time.Now()
	values := []interface{}{now.Format("02.01.2006 15:04:05"),
		emGoals.UsrGoals[0].UsrinvCurrentValuation,
		emGoals.UsrGoals[0].UsrinvCurrentValuationAssets,
		emGoals.UsrGoals[0].UsrinvTotalGains,
		emGoals.UsrGoals[0].UsrinvTotalDividendsReceived,
		emGoals.UsrGoals[0].UsrinvTotalContributions,
		emGoals.UsrGoals[0].UsrinvUnallocatedCreditAmount,
	}
	appendToOnvistSheet(srv, values)
}
