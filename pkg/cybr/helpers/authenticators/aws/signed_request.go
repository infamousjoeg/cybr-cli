package aws

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"time"
)

const (
	layoutISO = "2006-01-02"
	layoutUS  = "January 2, 2006"
	region    = "us-east-1"
	service   = "sts"
	host      = "sts.amazonaws.com"
)

func getCredentialScope(datestamp string) string {
	return datestamp + "/" + region + "/" + service + "/" + "aws4_request"
}

func hmacSHA256Hash(data string, key []byte) ([]byte, error) {
	h := hmac.New(sha256.New, key)
	_, err := h.Write([]byte(data))
	return h.Sum(nil), err
}

func getSignatureKey(key string, dateStamp, regionName, serviceName string) ([]byte, error) {
	kSecret := []byte("AWS4" + key)
	kDate, err := hmacSHA256Hash(dateStamp, kSecret)
	if err != nil {
		return nil, err
	}

	kRegion, err := hmacSHA256Hash(regionName, kDate)
	if err != nil {
		return nil, err
	}

	kService, err := hmacSHA256Hash(serviceName, kRegion)
	if err != nil {
		return nil, err
	}

	kSigning, err := hmacSHA256Hash("aws4_request", kService)
	if err != nil {
		return nil, err
	}

	return kSigning, nil
}

func signString(stringToSign string, signingKey []byte) (string, error) {
	hash, err := hmacSHA256Hash(stringToSign, signingKey)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash), nil
}

func getAuthorizationHeader(accessKey string, credentialScope, signedHeaders, signature string) string {
	algorithm := "AWS4-HMAC-SHA256"
	return fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s", algorithm, accessKey, credentialScope, signedHeaders, signature)
}

func headerAsJSONString(amzDate string, token, payloadHash, authorizationHeader string) string {
	headerTemplate := "{\"host\": \"%s\", \"x-amz-date\": \"%s\", \"x-amz-security-token\": \"%s\", \"x-amz-content-sha256\": \"%s\", \"authorization\": \"%s\"}"
	return fmt.Sprintf(headerTemplate, host, amzDate, token, payloadHash, authorizationHeader)
}

// GetDate returns the date in the format of yyyyMMdd
func GetDate(t time.Time) string {
	return t.UTC().Format("20060102")
}

// GetAmzDate returns the date in the format of yyyyMMdd'T'HHmmss'Z'
func GetAmzDate(t time.Time) string {
	return t.UTC().Format("20060102T150405Z")
}

// SHA256Hash returns the SHA256 hash of the input string
func SHA256Hash(input string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(input)))
}

// CreateStringToSign returns the string for AWS STS to sign
func CreateStringToSign(datestamp string, amzDate string, cannonicalRequest string) string {
	algorithm := "AWS4-HMAC-SHA256"
	credentialScope := getCredentialScope(datestamp)
	stringToSign := algorithm + "\n" + amzDate + "\n" + credentialScope + "\n" + SHA256Hash(cannonicalRequest)
	return stringToSign
}

// CreateCanonicalRequest returns the canonical request including signed headers from AWS STS
func CreateCanonicalRequest(amzdate string, token, signedHeaders, payloadHash string) string {
	canonicalURI := "/"
	canonicalQueryString := "Action=GetCallerIdentity&Version=2011-06-15"
	canonicalHeaders := "host:" + host + "\n" + "x-amz-content-sha256:" + payloadHash + "\n" + "x-amz-date:" + amzdate + "\n" + "x-amz-security-token:" + token + "\n"
	canonicalRequest := "GET" + "\n" + canonicalURI + "\n" + canonicalQueryString + "\n" + canonicalHeaders + "\n" + signedHeaders + "\n" + payloadHash
	return canonicalRequest
}

// GetAuthenticationRequest returns the signed request for AWS STS
func GetAuthenticationRequest(accessKey string, secretAccessKey string, token string, dateTime time.Time) (string, error) {
	amzDate := GetAmzDate(dateTime)
	dateStamp := GetDate(dateTime)

	signedHeaders := "host;x-amz-content-sha256;x-amz-date;x-amz-security-token"
	// payload is empty hence the hardcoded hash
	payloadHash := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	canonicalRequest := CreateCanonicalRequest(amzDate, token, signedHeaders, payloadHash)
	stringToSign := CreateStringToSign(dateStamp, amzDate, canonicalRequest)
	signingKey, err := getSignatureKey(secretAccessKey, dateStamp, region, service)
	if err != nil {
		return "", fmt.Errorf("Failed to get signature key. %s", err)
	}

	signature, err := signString(stringToSign, signingKey)
	if err != nil {
		return "", fmt.Errorf("Failed to sign string. %s", err)
	}

	authorizationHeader := getAuthorizationHeader(accessKey, getCredentialScope(dateStamp), signedHeaders, signature)
	return headerAsJSONString(amzDate, token, payloadHash, authorizationHeader), nil
}

// GetAuthenticationRequestNow returns the signed request for AWS STS using the current time
func GetAuthenticationRequestNow(accessKey string, secretAccessKey string, token string) (string, error) {
	return GetAuthenticationRequest(accessKey, secretAccessKey, token, time.Now())
}
