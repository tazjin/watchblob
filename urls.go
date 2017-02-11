package main

import "fmt"

const urlFormat string = "https://%s%s"
const triggerChallengeUri = "/?action=sslvpn_logon&fw_username=%s&fw_password=%s&style=fw_logon_progress.xsl&fw_logon_type=logon&fw_domain=Firebox-DB"
const responseUri = "/?action=sslvpn_logon&style=fw_logon_progress.xsl&fw_logon_type=response&response=%s&fw_logon_id=%d"

func templateChallengeTriggerUri(username *string, password *string) string {
	return fmt.Sprintf(triggerChallengeUri, *username, *password)
}

func templateResponseUri(logonId int, token *string) string {
	return fmt.Sprintf(responseUri, *token, logonId)
}

func templateUrl(baseUrl *string, uri string) string {
	return fmt.Sprintf("https://%s%s", *baseUrl, uri)
}
