/* This file has functions to interact with the conjur server */

package conjurapi

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func handleError(msg string, err error, exit bool) {
	if err != nil {
		fmt.Println(msg)
		fmt.Println(err)
		if exit {
			os.Exit(1)
		}
	}
}

// generate an authorization token to make requests to the conjur api
func GetToken(conjurHost, conjurAcct, hostName, apiKey string) string {
	url := fmt.Sprintf("https://%s/authn/%s/%s/authenticate", conjurHost, conjurAcct, hostName)

	//pull the apiKey from env variable, which will be set by the external secret operator in a kubernetes pod
	// apiKey := os.Getenv("apiKey")

	//create the request object
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(apiKey)))
	handleError("Error while creating the request object", err, true)

	//set headers
	req.Header.Set("Accept-Encoding", "base64")

	//create a transport object to skip ssl cert verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	//create the http client object
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	handleError("Error while making the http request", err, true)

	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	response := buf.String()

	return response
}

func GetSecret(conjurHost, conjurAcct, secretIdentifier, token string) string {
	authToken := fmt.Sprintf("Token token=\"%s\"", token)

	//double percent signs %% are used to escape the percent sign itself
	secret := secretIdentifier

	url := fmt.Sprintf("https://%s/secrets/%s/variable/%s", conjurHost, conjurAcct, secret)

	//create the request object
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
	handleError("Error while creating the request object", err, true)

	//set headers
	req.Header.Set("Authorization", authToken)
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept", "*/*")

	//create a transport object to skip ssl cert verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	//create the http client object
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	handleError("Error while making the http request", err, true)

	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	response := buf.String()

	if strings.Contains(response, "error") {
		return "invalid secret identifier"
	} else if strings.Contains(response, "Malformed") {
		return "malformed authorization token"
	}

	return response
}

func PullSecret(conjurHost, conjurAcct, hostname, secretIdentifier, apiKey string) string {
	hostnameParsed := strings.Replace(hostname, "/", "%2F", -1)
	secretIdentifierParsed := strings.Replace(secretIdentifier, "/", "%2F", -1)

	token := GetToken(conjurHost, conjurAcct, hostnameParsed, apiKey)

	password := GetSecret(conjurHost, conjurAcct, secretIdentifierParsed, token)

	return password
}
