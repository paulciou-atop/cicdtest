package configer

import (
	"nms/api/v1/devconfig"
	"time"

	"github.com/sirupsen/logrus"
)

// supportFields,
var supportFields = map[string][]string{
	"general": {"ip", "name"},
	"network": {"ip", "name", "mask", "gateway", "mac", "snmpEnable"},
	"snmp":    {"port", "public", "private", "ver"},
}

// contains check slice s contains e
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// allFail get a config which type is map[string]interface{} and return a string slice consist of
// all config field from the input config
func allFail(m map[string]interface{}) []string {
	var failFields = []string{}
	for field, _ := range m {
		failFields = append(failFields, field)
	}
	return failFields
}

// testUnsupport test input config payload has any un-supported fields?
// this function use global slice `supportFields`, configer service should implement thier own list
func testUnsupport(kind string, payload map[string]interface{}) []string {
	logrus.Infof("testUnsupport kind %s payload %v\n", kind, payload)
	var unsupportList = []string{}
	for field, _ := range payload {
		if !contains(supportFields[kind], field) {
			unsupportList = append(unsupportList, field)
		}
	}
	return unsupportList
}

func sessionSuccess(session *devconfig.SessionState) *devconfig.SessionState {
	session.State = "success"
	session.EndedTime = time.Now().String()
	// TODO: publish session finished
	return session
}

func sessionFail(session *devconfig.SessionState, msg string) *devconfig.SessionState {
	session.State = "fail"
	session.EndedTime = time.Now().String()
	session.Message = msg
	// TODO: publish session finished
	return session
}
