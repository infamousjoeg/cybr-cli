package identity

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/identity/requests"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/identity/responses"
)

// StartAuthentication starts the authentication process
func StartAuthentication(c api.Client, req requests.StartAuthentication) (*responses.StartAuthentication, error) {
	identityTenant := fmt.Sprintf("https://%s.id.cyberark.cloud", req.TenantID)
	url := fmt.Sprintf("%s/Security/StartAuthentication", identityTenant)
	s := fmt.Sprintf("{ \"TenantId\": \"%s\", \"User\": \"%s\", \"Version\": \"1.0\" }", req.TenantID, req.User)
	payload := strings.NewReader(s)

	client := &http.Client{
		Timeout: time.Second * 30, // Maximum of 30 secs
	}

	httpRequest, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return &responses.StartAuthentication{}, fmt.Errorf("Failed to start authentication. %s", err)
	}

	httpRequest.Header.Add("X-IDAP-NATIVE-CLIENT", "true")
	httpRequest.Header.Add("Content-Type", "application/json")

	dump, err := httputil.DumpRequestOut(httpRequest, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", dump)

	res, err := client.Do(httpRequest)
	if err != nil {
		fmt.Println(err)
		return &responses.StartAuthentication{}, fmt.Errorf("Failed to start authentication. %s", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return &responses.StartAuthentication{}, fmt.Errorf("Failed to start authentication. %s", err)
	}
	fmt.Println(string(body))

	StartAuthResponse := &responses.StartAuthentication{}
	err = json.Unmarshal(body, StartAuthResponse)
	return StartAuthResponse, err
}
