package tasks

import (
    "os"
    "time"
    "strings"
    "encoding/json"
    "ghostrunner/utils"
    "ghostrunner/logging"
)

func ProcessTasks(config *utils.Configuration, tasks []Task) {
    logging.Debug("tasks.controller", "ProcessTasks", "Starting to process tasks")

    for _, task := range tasks {
        logging.Debug("tasks.controller", "ProcessTasks", "Processing task")

        taskFolderLocation := config.ProcessingLocation + "/" + task.ExternalId

        logging.Debug("tasks.controller", "ProcessTasks", "Processing location set to " + taskFolderLocation)

        if _, err := os.Stat(taskFolderLocation); err == nil {
            logging.Debug("tasks.controller", "ProcessTasks", "Cleaning out previous processing location")

            os.RemoveAll(taskFolderLocation)
        }

        encryptionKey := config.ApiSecret[0:16] + strings.Replace(config.SessionId, "-", "", -1)[0:16]

        for _, encryptedScript := range task.EncryptedScripts {
            script := getTaskScript(utils.Decrypt(encryptionKey, encryptedScript))

            if len(script.Content) > 0 {
                var status, log string

                switch strings.ToLower(script.Type) {
                    case "node": 
                    script.Started = time.Now().UTC()
                    status, log = RunNodeScript(config, task.ExternalId, script.Content, script.Id)
                    case "commandline":
                    script.Started = time.Now().UTC()
                    status, log = RunCommandLineScript(config, task.ExternalId, script.Content, script.Id)
                }

                script.Completed = time.Now().UTC()
                script.Log = utils.Encrypt(encryptionKey, log)
                script.Status = status
            }

            task.Scripts = append(task.Scripts, script)
        }

        if _, err := os.Stat(taskFolderLocation); err == nil {
            logging.Debug("tasks.controller", "ProcessTasks", "Cleaning up processing location")

            os.RemoveAll(taskFolderLocation)
        }

        logging.Debug("tasks.controller", "ProcessTasks", "Finished processing task")

        UpdateTask(config.HostUrl, config.SessionId, task)
    }

    logging.Debug("tasks.controller", "ProcessTasks", "Finished processing tasks")
}

func getTaskScript(script []byte) (TaskScript) {
    taskScript := TaskScript{}
    json.Unmarshal(script, &taskScript)

    return taskScript
}