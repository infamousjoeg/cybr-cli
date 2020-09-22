package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	pasapi "github.com/infamousjoeg/pas-api-go/pkg/cybr/api"
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
	// Marshal struct data into JSON format for pretty print
	jsonVerify, err := json.Marshal(&resVerify)
	if err != nil {
		log.Fatalf("Unable to marshal struct to JSON for verify. %s", err)
	}
	fmt.Printf("Verify JSON:\r\n%s\r\n\r\n", jsonVerify)

	// Logon to PAS REST API Web Services
	token, errLogon := pasapi.Logon(hostname, username, password, authType, false)
	if errLogon != nil {
		log.Fatalf("Authentication failed. %s", errLogon)
	}
	fmt.Printf("Session Token:\r\n%s\r\n\r\n", token)

	// List Safes
	resListSafes := pasapi.ListSafes(hostname, token)
	jsonListSafes, err := json.Marshal(resListSafes)
	if err != nil {
		log.Fatalf("Unable to marshal struct to JSON for List Safes. %s", err)
	}
	fmt.Printf("List Safes JSON:\r\n%s\r\n\r\n", jsonListSafes)

	// List Safe Members
	resListSafeMembers := pasapi.ListSafeMembers(hostname, token, "D-Win-DomainAdmins")
	jsonListSafeMembers, err := json.Marshal(resListSafeMembers)
	if err != nil {
		log.Fatalf("Unable to marshal struct to JSON for List Safe Members. %s", err)
	}
	fmt.Printf("List Safe Members JSON:\r\n%s\r\n\r\n", jsonListSafeMembers)

	// List Applications
	resListApplications := pasapi.ListApplications(hostname, token, "\\", true)
	jsonListApplications, err := json.Marshal(resListApplications)
	if err != nil {
		log.Fatalf("Unable to marshal struct to JSON for List Applications. %s", err)
	}
	fmt.Printf("List Applications JSON:\r\n%s\r\n\r\n", jsonListApplications)

	// List Application Authentication Methods
	resListApplicationAuthenticationMethods := pasapi.ListApplicationAuthenticationMethods(hostname, token, "Ansible")
	jsonListApplicationAuthenticationMethods, err := json.Marshal(resListApplicationAuthenticationMethods)
	if err != nil {
		log.Fatalf("Unable to marshal struct to JSON for List Application Authentication Methods. %s", err)
	}
	fmt.Printf("List Application Authentication Methods JSON:\r\n%s\r\n\r\n", jsonListApplicationAuthenticationMethods)

	// Logoff PAS REST API Web Services
	success, errLogoff := pasapi.Logoff(hostname, token)
	if errLogoff != nil || success != true {
		log.Fatalf("Logoff failed. %s", errLogoff)
	}
	fmt.Println("Successfully logged off PAS REST API Web Services.")
}
