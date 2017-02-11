package main

import (
	"encoding/xml"
	"reflect"
	"testing"
)

func TestUnmarshalChallengeRespones(t *testing.T) {
	var testXml string = `
<?xml version="1.0" encoding="UTF-8"?>
<resp>
  <action>sslvpn_logon</action>
  <logon_status>4</logon_status>
  <auth-domain-list>
    <auth-domain>
      <name>RADIUS</name>
    </auth-domain>
  </auth-domain-list>
  <logon_id>441</logon_id>
  <chaStr>Enter Your 6 Digit Passcode </chaStr>
</resp>`

	var r Resp
	xml.Unmarshal([]byte(testXml), &r)

	expected := Resp{
		Action:      "sslvpn_logon",
		LogonStatus: 4,
		LogonId:     441,
		Challenge:   "Enter Your 6 Digit Passcode ",
	}

	assertEqual(t, expected, r)
}

func TestUnmarshalLoginError(t *testing.T) {
	var testXml string = `
<?xml version="1.0" encoding="UTF-8"?>
<resp>
  <action>sslvpn_logon</action>
  <logon_status>2</logon_status>
  <auth-domain-list>
    <auth-domain>
      <name>RADIUS</name>
    </auth-domain>
  </auth-domain-list>
  <errStr>501</errStr>
</resp>`

	var r Resp
	xml.Unmarshal([]byte(testXml), &r)

	expected := Resp{
		Action:      "sslvpn_logon",
		LogonStatus: 2,
		Error:       "501",
	}

	assertEqual(t, expected, r)
}

func TestUnmarshalLoginSuccess(t *testing.T) {
	var testXml string = `
<?xml version="1.0" encoding="UTF-8"?>
<resp>
  <action>sslvpn_logon</action>
  <logon_status>1</logon_status>
  <auth-domain-list>
    <auth-domain>
      <name>RADIUS</name>
    </auth-domain>
  </auth-domain-list>
</resp>
`
	var r Resp
	xml.Unmarshal([]byte(testXml), &r)

	expected := Resp{
		Action:      "sslvpn_logon",
		LogonStatus: 1,
	}

	assertEqual(t, expected, r)
}

func assertEqual(t *testing.T, expected interface{}, result interface{}) {
	if !reflect.DeepEqual(expected, result) {
		t.Errorf(
			"Unmarshaled values did not match.\nExpected: %v\nResult: %v\n",
			expected, result,
		)

		t.Fail()
	}
}
