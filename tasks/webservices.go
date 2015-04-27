package tasks

import (
    "bytes"
    "strconv"
    "net/http"
    _ "crypto/sha512"
    "io/ioutil"
    "encoding/json"
    "errors"
    "ghostrunner/utils"
    "ghostrunner/logging"
)

func GetAll(hostUrl, taskServerSessionId string) ([]Task) {   
    var tasks []Task

    logging.Debug("tasks.webservices", "GetAll", "Getting tasks from '" + hostUrl + "/tasks'")

    request, err := http.NewRequest("GET", hostUrl + "/1_0/tasks", nil)
    request.Header.Add("sessionid", taskServerSessionId)

    client := &http.Client{}
    response, err := client.Do(request)
    
    if err != nil {
        logging.Error("tasks.webservices", "GetAll", "Error performing request", err)
    } else {
        defer response.Body.Close()

        logging.Debug("tasks.webservices", "GetAll", "Parsing out response body")

        contents, err := ioutil.ReadAll(response.Body)

        if err != nil {
            logging.Error("tasks.webservices", "GetAll", "Error reading response", err)
        }

        var objmap map[string]*json.RawMessage
        json.Unmarshal(contents, &objmap)

        logging.Debug("tasks.webservices", "GetAll", "Checking response code")

        if _, ok := objmap["code"]; ok {
            var code int
            json.Unmarshal(*objmap["code"], &code)

            if code == 200 {
                logging.Debug("tasks.webservices", "GetAll", "Successfully retrieved tasks")

                json.Unmarshal(*objmap["tasks"], &tasks)
            } else if code == 500 {
                if _, ok := objmap["message"]; ok {
                    var message string
                    json.Unmarshal(*objmap["message"], &message)

                    logging.Error("tasks.webservices", "GetAll", "(500) Unable to retrieve all server tasks", errors.New(message))
                }
            } else {
                logging.Debug("tasks.webservices", "GetAll", "(" + strconv.Itoa(code) + ") Unable to retrieve all server tasks")

            }
        }
    }

    return tasks
}

func UpdateTask(hostUrl, taskServerSessionId string, task Task) {
    logging.Debug("tasks.webservices", "Update", "Updating task at '" + hostUrl + "/tasks'")

    logJSON := "{ \"scripts\": ["

    for i := 0; i < len(task.Scripts); i++ {
        if i != 0 {
            logJSON += ","            
        }

        logJSON += "{ \"id\":\"" + strconv.Itoa(task.Scripts[i].Id) + "\", \"status\":\"" + task.Scripts[i].Status + "\", \"log\":\"" + utils.EscapeForJSON(task.Scripts[i].Log) + "\", \"started\":\"" + task.Scripts[i].Started.String() + "\", \"completed\":\"" + task.Scripts[i].Completed.String() + "\" }"
    }

    logJSON += "] }"

    request, _ := http.NewRequest("PUT", hostUrl + "/1_0/tasks/" + task.ExternalId, bytes.NewBuffer([]byte(logJSON)))
    request.Header.Add("Content-Type", "application/json")
    request.Header.Add("sessionid", taskServerSessionId)

    client := &http.Client{}
    response, err := client.Do(request)

    if err != nil {
        logging.Error("tasks.webservices", "UpdateTask", "Unable to perform web request", err)
    } else {
        defer response.Body.Close()

        contents, err := ioutil.ReadAll(response.Body)

        if err != nil {
            logging.Error("tasks.webservices", "UpdateTask", "Error reading response", err)
        }

        var objmap map[string]*json.RawMessage
        json.Unmarshal(contents, &objmap)

        logging.Debug("tasks.webservices", "UpdateTask", "Checking response code")

        if _, ok := objmap["code"]; ok {
            var code int
            json.Unmarshal(*objmap["code"], &code)

            if code == 500 {
                if _, ok := objmap["message"]; ok {
                    var message string
                    json.Unmarshal(*objmap["message"], &message)

                    logging.Error("tasks.webservices", "UpdateTask", "(500) Unable to authenticate task runner", errors.New(message))
                }
            } else {
                logging.Debug("tasks.webservices", "UpdateTask", "(" + strconv.Itoa(code) + ") Unable to authenticate task runner")
            }
        }
    }
}