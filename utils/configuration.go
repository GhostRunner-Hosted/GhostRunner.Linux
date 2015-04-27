package utils

import (
    "os"
    "encoding/json"
    "ghostrunner/logging"
    "errors"
)

type Configuration struct {
	HostUrl string
	ProcessingLocation string
    NodeLocation string
    NpmLocation	string
    RunnerId string
    ApiKey string
    ApiSecret string
	SessionId string
}

func LoadConfiguration() (Configuration, error) {
	if _, err := os.Stat("/etc/ghostrunner.conf"); os.IsNotExist(err) {
        logging.Error("utils.configuration", "LoadConfiguration", "Unable to fing configuration file '/etc//ghostrunner.conf'", err)

		return Configuration{}, errors.New("Unable to fing configuration file '/etc/ghostrunner.conf'")
	}
	
	file, _ := os.Open("/etc/ghostrunner.conf")

	decoder := json.NewDecoder(file)

	configuration := Configuration{}

	err := decoder.Decode(&configuration)
    
    if err != nil {
        logging.Error("utils.configuration", "LoadConfiguration", "Error decoding configuration file", err)
        
        return Configuration{}, errors.New("Error decoding configuration file '/etc/ghostrunner.conf'")
    }

    if len(configuration.RunnerId) == 0 {
        runnerId, _ := GenerateUUID()

        configuration.RunnerId = runnerId

        UpdateConfiguration(&configuration) 
    }

	return configuration, nil
}

func UpdateConfiguration(config *Configuration) {
    if _, err := os.Stat("/etc/ghostrunner.conf"); err == nil {
        err := os.Remove("/etc/ghostrunner.conf")

        if err != nil {
            logging.Error("utils.configuration", "UpdateConfiguration", "Error deleting previous configuration", err)

            return
        }
    }

    configFile, err := os.Create("/etc/ghostrunner.conf")

    if err != nil {
        logging.Error("utils.configuration", "UpdateConfiguration", "Error creating configuration", err)

        return
    }
    
    defer configFile.Close()

	stringifiedConfig, _ := json.Marshal(config)

    _, err = configFile.Write([]byte(stringifiedConfig))
}
