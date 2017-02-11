package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// The XML response returned by the WatchGuard server
type Resp struct {
	Action      string `xml:"action"`
	LogonStatus int    `xml:"logon_status"`
	LogonId     int    `xml:"logon_id"`
	Error       string `xml:"errStr"`
	Challenge   string `xml:"chaStr"`
}

func main() {
	args := os.Args[1:]

	if len(args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: watchblob <vpn-host> <username> <password>\n")
		os.Exit(1)
	}

	host := args[0]
	username := args[1]
	password := args[2]

	challenge, err := triggerChallengeResponse(&host, &username, &password)

	if err != nil || challenge.LogonStatus != 4 {
		fmt.Fprintln(os.Stderr, "Did not receive challenge from server")
		fmt.Fprintf(os.Stderr, "Response: %v\nError: %v\n", challenge, err)
		os.Exit(1)
	}

	token := getToken(&challenge)
	err = logon(&host, &challenge, &token)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Logon failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Login succeeded, you may now (quickly) authenticate OpenVPN with %d as your password\n", token)
}

func triggerChallengeResponse(host *string, username *string, password *string) (r Resp, err error) {
	return request(templateUrl(host, templateChallengeTriggerUri(username, password)))
}

func getToken(challenge *Resp) string {
	fmt.Println(challenge.Challenge)

	reader := bufio.NewReader(os.Stdin)
	token, _ := reader.ReadString('\n')

	return strings.TrimSpace(token)
}

func logon(host *string, challenge *Resp, token *string) (err error) {
	resp, err := request(templateUrl(host, templateResponseUri(challenge.LogonId, token)))
	if err != nil {
		return
	}

	if resp.LogonStatus != 1 {
		err = fmt.Errorf("Challenge/response authentication failed: %v", resp)
	}

	return
}

func request(url string) (r Resp, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	decoder := xml.NewDecoder(resp.Body)

	err = decoder.Decode(&r)
	return
}
