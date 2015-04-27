package logging

import (
    "os"
    "encoding/json"
)

type Configuration struct {
	ErrorLogLocation string
	ErrorMaxFileSizeKB int
	ErrorRollover bool
	ErrorMaxRollerOvers int
	DebugEnabled bool
	DebugLogLocation string
	DebugMaxFileSizeKB int
	DebugRollover bool
	DebugMaxRollerOvers int
}

func LoadConfiguration() (Configuration) {
	if _, err := os.Stat("/etc/ghostrunner.log.conf"); os.IsNotExist(err) {
		return createDefaultConfig()
	}
	
	file, _ := os.Open("/etc/ghostrunner.log.conf")

	decoder := json.NewDecoder(file)

	configuration := createDefaultConfig()

	err := decoder.Decode(&configuration)

	if err != nil {
		return createDefaultConfig()
	}

	return configuration
}

func createDefaultConfig() (Configuration) {
	configuration := Configuration{}

	configuration.ErrorLogLocation = "/var/log/ghostrunner/error.log"
	configuration.ErrorMaxFileSizeKB = 10
	configuration.DebugLogLocation =  "/var/log/ghostrunner/debug.log"
	configuration.DebugEnabled = false
	configuration.DebugMaxFileSizeKB = 10

	return configuration
}
