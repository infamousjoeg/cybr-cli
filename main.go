package main

import (
	"encoding/json"
	"fmt"
	"log"

	pasapi "github.com/infamousjoeg/pas-api-go/pkg/api"
)

const (
	hostname = "https://cyberark.joegarcia.dev"
)

func main() {
	// Verify PAS REST API Web Services
	resVerify, errVerify := pasapi.ServerVerify(hostname)
	if errVerify != nil {
		log.Fatalf("Verification failed. %s", errVerify)
	}
	b, err := json.Marshal(resVerify)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))

	// Authenticate with PAS REST API Web Services
	token, errAuth := pasapi.Authenticate(hostname, "username", "password", "ldap")
	if errAuth != nil {
		log.Fatalf("Authentication failed. %s", errAuth)
	}
	fmt.Println(token)
}
