package identity

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/httpjson"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/identity/requests"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/identity/responses"
)

func checkURL(client api.Client, req requests.AdvanceAuthentication, c chan string) {
	for {
		identityTenant := fmt.Sprintf("https://%s.id.cyberark.cloud", client.TenantID)
		url := fmt.Sprintf("%s/Security/AdvanceAuthentication", identityTenant)

		headers := http.Header{}
		headers.Add("X-IDAP-NATIVE-CLIENT", "true")
		headers.Add("Content-Type", "application/json")

		res, err := httpjson.SendRequestRawWithHeaders(url, "POST", headers, req, client.InsecureTLS, client.Logger)
		if err != nil {
			log.Fatalf("Failed to reach Identity to check OOB status. %s", err)
		}
		advanceAuthResponse := &responses.Authentication{}
		err = json.Unmarshal(res, advanceAuthResponse)
		if err != nil || advanceAuthResponse.Result.Summary == "OobPending" {
			time.Sleep(1 * time.Second)
			continue
		}
		c <- advanceAuthResponse.Result.Token
	}
}

func waitForInput(c chan string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the one-time passcode or click the link: ")
	for scanner.Scan() {
		c <- scanner.Text()
	}
}

// GetOOBPending will return the one-time passcode or successful link click
func GetOOBPending(c api.Client, req requests.AdvanceAuthentication) (string, error) {
	responseChannel := make(chan string)

	go checkURL(c, req, responseChannel)
	go waitForInput(responseChannel)

	return <-responseChannel, nil
}
