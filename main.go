package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	pasapi "github.com/infamousjoeg/pas-api-go/pkg/api"
)

var (
	hostname = os.Getenv("PAS_HOSTNAME")
	username = os.Getenv("PAS_USERNAME")
	password = os.Getenv("PAS_PASSWORD")
	authType = os.Getenv("PAS_AUTH_TYPE")
)

func main() {
	// Verify PAS REST API Web Services
	resVerify := pasapi.ServerVerify(hostname)
	jsonVerify, err := json.Marshal(&resVerify)
	if err != nil {
		log.Fatalf("Unable to marshall struct to JSON for verify. %s", err)
	}
	fmt.Printf("Verify JSON:\r\n%s\r\n\r\n", jsonVerify)

	// Logon to PAS REST API Web Services
	token, errLogon := pasapi.Logon(hostname, username, password, authType, false)
	if errLogon != nil {
		log.Fatalf("Authentication failed. %s", errLogon)
	}
	fmt.Printf("Session Token:\r\n%s\r\n\r\n", token)

	// Logoff PAS REST API Web Services
	success, errLogoff := pasapi.Logoff(hostname, token)
	if errLogoff != nil || success != true {
		log.Fatalf("Logoff failed. %s", errLogoff)
	}
	fmt.Println("Successfully logged off PAS REST API Web Services.")
}
