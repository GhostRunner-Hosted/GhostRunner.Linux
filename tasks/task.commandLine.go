package tasks

import (
    "strings"
    "strconv"
    "runtime"
    "bytes"
    "bufio"
    "io"
    "time"
    "os"
    "os/exec"
    "ghostrunner/utils"
    "ghostrunner/logging"
)

const (
    gitTimeout = 10
)

type GitRepository struct {
    Location string `json:"location"`
    Username string `json:"username"`
    Password string `json:"password"`
}

func RunCommandLineScript(config *utils.Configuration, taskId, script string, scriptId int) (string, string) {
    logging.Debug("task.commandLine", "RunCommandLineScript", "Starting to process commandLine task");
    
    if strings.ToLower(runtime.GOOS) == "linux" {
        scriptLog := ""
        status := "Completed"

        var shellScriptName string
        var shellFileLocation string

        logging.Debug("task.commandLine", "RunCommandLineScript", "System detected as Linux")

        shellScriptName = taskId + "_" + strconv.Itoa(scriptId) + ".sh"

        script = strings.TrimSpace(script)

        if !strings.HasPrefix(strings.ToLower(script), "#!/bin/bash") {
            script = "#!/bin/bash\n" + script
        }
    
        if len(shellScriptName) > 0 {
            taskFolderLocation := config.ProcessingLocation + "/" + taskId

            if _, err := os.Stat(taskFolderLocation); os.IsNotExist(err) {
                logging.Debug("task.commandLine", "RunCommandLineScript", "Creating task processing location")
                err := os.Mkdir(taskFolderLocation, 0777)

                if err != nil {
                    logging.Error("task.commandLine", "RunCommandLineScript", "Error creating task processing location " + taskFolderLocation, err)
                }
            }

            shellFileLocation = config.ProcessingLocation + "/" + taskId + "/" + shellScriptName

            if _, err := os.Stat(shellFileLocation); err == nil {
                err := os.Remove(shellFileLocation)

                if err != nil {
                    logging.Error("task.commandLine", "RunCommandLineScript", "Error deleting previous shell script at " + shellFileLocation, err)

                    return "Errored", "Error deleting previous shell script"
                }
            }

            shellFile, err := os.Create(shellFileLocation)

            if err != nil {
                logging.Error("task.commandLine", "RunCommandLineScript", "Error creating shell script at " + shellFileLocation, err)
                return "Errored", "Error creating shell script"
            }
            
            defer shellFile.Close()
            defer os.Remove(shellFileLocation)

            _, err = shellFile.Write([]byte(script))

            if err != nil {
                logging.Error("task.commandLine", "RunCommandLineScript", "Error writing shell script to " + shellFileLocation, err)
                return "Errored", "Error writing shell script"
            }

            cmd := exec.Command("sh", shellScriptName)
            cmd.Dir = taskFolderLocation

            stdout, err := cmd.StdoutPipe()
            if err != nil {
                logging.Error("task.commandLine", "RunCommandLineScript", "Error creating command line pipe", err)
            }

            go func() {
                time.Sleep(nodeTimeout * time.Minute)
                cmd.Process.Kill()

                logging.Error("task.commandLine", "RunCommandLineScript", "Command line script timed out", err)
                
                status = "Errored"
            }()

            go func() {
                err = cmd.Run()
                if err != nil {
                    logging.Error("task.commandLine", "RunCommandLineScript", "Error running command line script", err)

                    status = "Errored"
                }
            }()

            var buffer bytes.Buffer
            writer := bufio.NewWriter(&buffer)
            defer writer.Flush()

            io.Copy(writer, stdout)

            scriptLog += buffer.String()

            return status, scriptLog
        } else {
            return "Errored", ""
        }
    }

    return "Errored", "Incorrect operating system"
}
