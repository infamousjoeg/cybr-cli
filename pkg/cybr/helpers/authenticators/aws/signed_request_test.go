package aws_test

import (
	"testing"
	"time"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/authenticators/aws"
)

const (
	layoutISO                     = "2006-01-02T15:04:05-0700"
	layoutUS                      = "January 2, 2006"
	date                          = "2006-01-02T15:04:05-0700"
	expectedDateAsString          = "19700101"
	expectedAmzDateAsString       = "19700101T010101Z"
	expectedCanonicalRequest      = "GET\n/\nAction=GetCallerIdentity&Version=2011-06-15\nhost:sts.amazonaws.com\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\nx-amz-date:20060102T220405Z\nx-amz-security-token:thisIsMyToken\n\nhost;x-amz-content-sha256;x-amz-date;x-amz-security-token\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	expectedCreateStringToSign    = "AWS4-HMAC-SHA256\n" + expectedAmzDateAsString + "\n" + expectedDateAsString + "/us-east-1/sts/aws4_request\n6e861cade935d3d9bad0fef193202898f8dc45e33941d763a6387a5a77d077fa"
	sessionToken                  = "thisIsMyToken"
	signedHeaders                 = "host;x-amz-content-sha256;x-amz-date;x-amz-security-token"
	payloadHash                   = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	expectedAuthenticationRequest = "{\"host\": \"sts.amazonaws.com\", \"x-amz-date\": \"19700101T060101Z\", \"x-amz-security-token\": \"sessionToken\", \"x-amz-content-sha256\": \"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\", \"authorization\": \"AWS4-HMAC-SHA256 Credential=accessKey/19700101/us-east-1/sts/aws4_request, SignedHeaders=host;x-amz-content-sha256;x-amz-date;x-amz-security-token, Signature=9c9d03c90a55b460a7bb4d3caf4a2900518c7d51a4ad7e8f5aeadbd3050edaed\"}"
)

func assertStringEquals(t *testing.T, actual string, expected string) {
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestGetDate(t *testing.T) {
	expected := "20060102"
	time, _ := time.Parse(layoutISO, date)
	actual := aws.GetDate(time)
	assertStringEquals(t, actual, expected)
}

func TestGetAmzDate(t *testing.T) {
	expected := "20060102T220405Z"
	time, _ := time.Parse(layoutISO, date)
	actual := aws.GetAmzDate(time)
	assertStringEquals(t, actual, expected)
}

func TestSHA256Hash(t *testing.T) {
	expected := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
	actual := aws.SHA256Hash("hello world")
	assertStringEquals(t, actual, expected)
}

func TestCreateStringToSign(t *testing.T) {
	actual := aws.CreateStringToSign(expectedDateAsString, expectedAmzDateAsString, expectedCanonicalRequest)
	assertStringEquals(t, actual, expectedCreateStringToSign)
}

func TestCreateCanonicalRequest(t *testing.T) {
	expected := "20060102T220405Z"
	time, _ := time.Parse(layoutISO, date)
	amzDate := aws.GetAmzDate(time)
	assertStringEquals(t, amzDate, expected)

	actual := aws.CreateCanonicalRequest(amzDate, sessionToken, signedHeaders, payloadHash)
	assertStringEquals(t, actual, expectedCanonicalRequest)
}

func TestGetAuthenticationRequest(t *testing.T) {
	time, _ := time.Parse(layoutISO, "1970-01-01T06:01:01-0000")
	actual, err := aws.GetAuthenticationRequest("accessKey", "secretKey", "sessionToken", time)
	if err != nil {
		t.Errorf("Failed to get authentication request. %s", err)
	}
	assertStringEquals(t, actual, expectedAuthenticationRequest)
}

// func TestGetAuthenticationRequestNow(t *testing.T) {
// 	accessKey := "ASIAZB7BKVSO5BQRVRP7"
// 	secretKey := "muF+nDmEMaZtC1mG/9Y+/4Irq44Tv3kUXmksyact"
// 	sessionToken := "IQoJb3JpZ2luX2VjEBIaCXVzLWVhc3QtMSJHMEUCIQCStfOCGFNPL2lidIhSFfi+3n+B0pKle3zze+VO4WhL3AIgLwJyHPaKP1KLLjoRyvnfh7nviKQGrb5EN2gArLk27FkqtAMIOhABGgw2MjI3MDU5NDU3NTciDIXDygKJq+O/FMRTcCqRA30hcRLo8uQdEzQCv3fMH+Z7ZkgaRpPpQLXh3QpKKRXiflF8JkE21Jtt+jOoXuV9n6IElvGq4mHsMYX4tFzxl8FR+EEcjO3+z3L74dH+ByCv4jlJQaoaHOWOznpR6NKxKsU5VVttD2Rd8ZxQndKyvbB+6TP2NhkuyggQJFHqDEzbIeoFsWvHI0AGF2B8UkjHp1qTiBQ/sltnuemOMoblpK9GqumEq0UD8iDAb9durkIS9A0Mo7XQjd5BlddImNTRVDwW7LW+z/LFZLoeVdwzJIQWBGjkkgFOhSkz0sP5qr/iMPOfEYjvfTDXlY+2okWZHbTS3zx2yR/9hB4tZCHSBcRkWkRpjBqBHgveeekeaKRiTtOXBwAWTKP8CyxKUueG5fgdTOAxFqMWkJZxfPsUGu5Rx02xkrdnmm8LGZ/uDlzjdg/fgX5B0bsOS+OKkCcqaap8hGArOJfNAmKOP07tX/FS/6v1GYo1pNyQl0q4sQfbF52csbqYq2l7XMOg3cDU9Mmx2vbA9A1ADKYkNy/M0F+vMLeXyvsFOusBbagVAza217PYNOEnvxKXjUrHSwZaqlYISgp8QO9PSW+nCMeBDNl/JwVsC1cKaQ9tuof79qQtGfuEHW2I+MDjss8bQULn2pyTgswUllVfSdIDkWR9qUyQrfT5PyDpQrwfHyqcUwk9XBivpfw7/hTy7qlIL46zevXJwg/X6Yq1pkvteaNzHuV9TXUqPokdF/Ra/E7tWHLaMmUHO8kLWOMShxcgYPlXSgzt6q11dST+7ya+s6ffPOOz5yNOo3NRjUfIH2zQfvSZfvgqoyVYPgba3mnRhtCRnOc0926DLSi3lhLekZDcO4bjNcQbYQ=="
// 	actual, _ := utils.GetAuthenticationRequestNow(accessKey, secretKey, sessionToken)
// 	t.Errorf("Get authentication request now. %s", actual)
// }
