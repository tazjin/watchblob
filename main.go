package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"net/http"
	"os"
	"strings"
	"syscall"
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

	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Usage: watchblob <vpn-host>")
		os.Exit(1)
	}

	host := args[0]

	username, password, err := readCredentials()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read credentials: %v\n", err)
	}

	fmt.Printf("Requesting challenge from %s as user %s\n", host, username)
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

	fmt.Printf("Login succeeded, you may now (quickly) authenticate OpenVPN with %s as your password\n", token)
}

func readCredentials() (string, string, error) {
	fmt.Printf("Username: ")
	reader := bufio.NewReader(os.Stdin)
	username, err := reader.ReadString('\n')

	fmt.Printf("Password: ")
	password, err := terminal.ReadPassword(syscall.Stdin)

	// If an error occured, I don't care about which one it is.
	return strings.TrimSpace(username), strings.TrimSpace(string(password)), err
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
