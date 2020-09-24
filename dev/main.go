package main

import (
	"fmt"
	"log"
	"os"

	pasapi "github.com/infamousjoeg/pas-api-go/pkg/cybr/api"
)

var (
	hostname = os.Getenv("PAS_HOSTNAME")
	username = os.Getenv("PAS_USERNAME")
	password = os.Getenv("PAS_PASSWORD")
)

func main() {
	client := pasapi.Client{
		Hostname: hostname,
		AuthType: "cyberark",
	}

	response, err := client.ServerVerify()
	if err != nil {
		log.Fatalf("Failed to get ServerVerify information for the PVWA. %s", err)
		return
	}
	fmt.Printf("Server ID: %s\nServer name: %s\n", response.ServerID, response.ServerName)

	credentials := pasapi.LogonRequest{
		Username: username,
		Password: password,
	}

	err = client.Logon(credentials)
	if err != nil {
		log.Fatalf("Failed to Logon to the PVWA. %s", err)
		return
	}

	safes, err := client.ListSafes()
	if err != nil {
		log.Fatalf("%s", err)
		return
	}

	// Lets iterate over each safe
	for _, s := range safes.Safes {
		// Get the members of each safe
		members, err := client.ListSafeMembers(s.SafeName)
		if err != nil {
			log.Fatalf("Failed to list members of safe '%s'. %s", s.SafeName, err)
			return
		}

		// Iterate each member in this safe and print out safe and members
		fmt.Printf("%s members\n", s.SafeName)
		for _, m := range members.Members {
			fmt.Printf("\t- %s\n", m.Username)
		}
	}

	apps, err := client.ListApplications("\\")
	if err != nil {
		log.Fatalf("%s", err)
		return
	}

	// Iterate through the applications
	for _, a := range apps.Application {
		// Get authentication methods for each appliucation
		authMethods, err := client.ListApplicationAuthenticationMethods(a.AppID)
		if err != nil {
			log.Fatalf("Failed to list auth methods for application '%s'. %s", a.AppID, err)
			return
		}

		// Print out app ID and authentication method types and values
		fmt.Printf("%s authentication methods\n", a.AppID)
		for _, m := range authMethods.Authentication {
			fmt.Printf("\t- %s = %s\n", m.AuthType, m.AuthValue)
		}
	}

	err = client.Logoff()
	if err != nil {
		log.Fatalf("Failed to log off. %s", err)
		return
	}
}
