package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"

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
	var amAll emAll
	_, _, errs = request.
		SetDebug(false).
		Get("https://gateway.etfmatic.com/user/1.0/retrieve").
		Param("session_token", emSignin.Token).
		Param("system_identifier", emSignin.ClientSystems.EuLive).
		EndStruct(&amAll)

	if errs != nil {
		fmt.Println(errs)
	}
	spew.Dump(amAll.UserAccountsSummary.TotalContributions)
}
