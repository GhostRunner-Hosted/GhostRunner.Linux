package tasks

import (
    "strings"
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
    npmTimeout = 10
)
    
func RunNpmScript(config *utils.Configuration, taskId, packageName string) (string) {
    logging.Debug("task.npm", "RunNpmScript", "Starting to process spm task");
    
    logging.Debug("task.npm", "RunNpmScript", "Checking operating system as this allows us to check for node");

    if strings.ToLower(runtime.GOOS) == "linux" {
        if _, err := os.Stat(config.NodeLocation); os.IsNotExist(err) {
            return "node.js is not installed"
        }

        logging.Debug("task.npm", "RunNpmScript", "System detected as Linux");

        if _, err := os.Stat(config.NpmLocation); os.IsNotExist(err) {
            return "npm is not installed"
        }

        taskFolderLocation := config.ProcessingLocation + "/" + taskId

        logging.Debug("task.npm", "RunNpmScript", "Checking for task folder location");

        if _, err := os.Stat(taskFolderLocation); os.IsNotExist(err) {
            logging.Debug("task.npm", "RunNpmScript", "Creating task folder  location");
            os.Mkdir(taskFolderLocation, 0777)
        }

        cmd := exec.Command("npm", "install", packageName)
        cmd.Dir = taskFolderLocation

        stdout, err := cmd.StdoutPipe()
        if err != nil {
            logging.Error("task.npm", "RunNpmScript", "Error creating npm pipe " + packageName, err)
        }

        go func() {
            time.Sleep(npmTimeout * time.Minute)
            cmd.Process.Kill()

            logging.Error("task.npm", "RunNpmScript", "npm install timed out " + packageName, err)
        }()

        go func() {
            err = cmd.Run()

            if err != nil {
                logging.Error("task.npm", "RunNpmScript", "Error retrieving npm package " + packageName, err)
            }
        }()

        var buffer bytes.Buffer
        writer := bufio.NewWriter(&buffer)
        defer writer.Flush()

        io.Copy(writer, stdout)

        return buffer.String()
    }
    
    return "Incorrect operating system"
}
