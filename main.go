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
	resVerify, errVerify := pasapi.ServerVerify(hostname)
	if errVerify != nil {
		log.Fatalf("Verification failed. %s", errVerify)
	}

	// Marshal (convert) returned map string interface to JSON
	b, err := json.Marshal(resVerify)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Verify JSON:\r\n%s\r\n\r\n", string(b))

	// Authenticate with PAS REST API Web Services
	token, errAuth := pasapi.Authenticate(hostname, username, password, authType)
	if errAuth != nil {
		log.Fatalf("Authentication failed. %s", errAuth)
	}
	fmt.Printf("Session Token:\r\n%s\r\n", token)
}
