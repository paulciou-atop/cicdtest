package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func Result(status int, message string, response any) map[string]any {
	// create new map and set key value
	m := make(map[string]any)
	// status
	m["status"] = status
	// message
	m["message"] = message
	// response
	m["response"] = response
	// return map
	return m
}

// normalize port
func NormalizePort(port string) string {
	// Ex: 134 > :134 | :8080" -> :8080 | abc -> abc
	var rePort = regexp.MustCompile("(?m)^([1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$")
	var rePass = regexp.MustCompile("(?m)^:([1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$")
	if rePass.MatchString(port) {
		return port
	}
	if rePort.MatchString(port) {
		return ":" + port
	}
	return port
}

// update session state
func UpdateSessionState(oldstates string, newstate string, key string) string {
	// create init map
	data := map[string]string{}
	// update new state
	states := strings.Split(oldstates, "|")
	for _, i := range states {
		kv := strings.Split(i, ":")
		data[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}
	data[key] = newstate
	// combine new string
	return fmt.Sprintf("gwd:%s|snmp:%s", data["gwd"], data["snmp"])
}
