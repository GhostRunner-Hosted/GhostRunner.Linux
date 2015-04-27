package taskRunner

import (
	"ghostrunner/tasks"
    "ghostrunner/utils"
    "ghostrunner/logging"
)

func ValidateRunner(config *utils.Configuration) (int) {
  	returnCode := CheckSessionId(config.HostUrl, config.SessionId)

  	return returnCode
}

func Authenticate(config *utils.Configuration) (string) {
	logging.Debug("taskrunner.controller", "Authenticate", "Connecting to the Sqlite database")
	
	runner := AuthenticateTaskRunner(config.HostUrl, config.RunnerId, config.ApiKey)

	return runner.SessionId
}

func GetTasks(config *utils.Configuration) ([]tasks.Task) {
	logging.Debug("taskrunner.controller", "GetTasks", "Connecting to the Sqlite database")

	var unprocessedTasks []tasks.Task

	if (len(config.SessionId) > 0) {
		unprocessedTasks = tasks.GetAll(config.HostUrl, config.SessionId)
	}

	return unprocessedTasks
}