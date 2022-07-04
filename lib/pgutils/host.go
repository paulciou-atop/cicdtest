package pgutils

import (
	"os"
)

//default value

const (
	host = ""
)

func isRunningInDockerContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	return false
}

// getDBHost
// TODO -- add config
func getDBHost() string {
	if isRunningInDockerContainer() {
		return "postgresql:5432"
		//return fmt.Sprintf("postgres://%s:%s@postgresql:5432/nms", user, pass)
	}
	return "localhost:5432"
	// return fmt.Sprintf("postgres://%s:%s@localhost:5432/nms", user, pass)

}
