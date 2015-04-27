package taskRunner

import (
    "strconv"
    "os"
	"bytes"
    "net/http"
    _ "crypto/sha512"
    "io/ioutil"
    "encoding/json"
    "errors"
    "ghostrunner/logging"
)

func CheckSessionId(hostUrl, sessionId string) (int) {
    req, _ := http.NewRequest("GET", hostUrl + "/1_0/runner/validate/" + sessionId, nil)

    client := &http.Client{}
    response, err := client.Do(req)

    if err != nil {
        logging.Error("taskRunner.webservices", "ValidateRunner", "Unable to perform web request", err)

        return 500
    } else {
        defer response.Body.Close()

        contents, err := ioutil.ReadAll(response.Body)

        if err != nil {
            logging.Error("taskRunner.webservices", "ValidateRunner", "Error reading response", err)

            return 500
        } else {
            var objmap map[string]*json.RawMessage
            json.Unmarshal(contents, &objmap)

            logging.Debug("taskRunner.webservices", "ValidateRunner", "Checking response code")

            if _, ok := objmap["code"]; ok {
                var code int
                json.Unmarshal(*objmap["code"], &code)

                return code
            }
        }
    }

    return 500
}

func AuthenticateTaskRunner(hostUrl, runnerId, apiKey string) (TaskRunner) {
	osHostName, _ := os.Hostname()

    logging.Debug("taskRunner.webservices", "Authenticate", "Authenticating task runner with '" + hostUrl + "/runner/authenticate'")

    var jsonStr = []byte("{ \"runnerId\":\"" + runnerId + "\", \"machineName\":\"" + osHostName + "\", \"apiKey\":\"" + apiKey + "\" }")

    req, _ := http.NewRequest("POST", hostUrl + "/1_0/runner/authenticate", bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    response, err := client.Do(req)

    if err != nil {
        logging.Error("taskRunner.webservices", "Authenticate", "Unable to perform web request '" + hostUrl + "/1_0/runner/authenticate'", err)

        return TaskRunner{}
    } else {
        defer response.Body.Close()

        contents, err := ioutil.ReadAll(response.Body)

        if err != nil {
            logging.Error("taskRunner.webservices", "Authenticate", "Error reading response", err)
        }

        var objmap map[string]*json.RawMessage
        json.Unmarshal(contents, &objmap)

        logging.Debug("taskRunner.webservices", "Authenticate", "Checking response code")

        if _, ok := objmap["code"]; ok {
            var code int
            json.Unmarshal(*objmap["code"], &code)

            if code == 200 {
                logging.Debug("taskRunner.webservices", "Authenticate", "Successfully authenticated task runner")

                var taskRunner TaskRunner
                json.Unmarshal(*objmap["taskServer"], &taskRunner)

            	return taskRunner
            } else if code == 500 {
                if _, ok := objmap["message"]; ok {
                    var message string
                    json.Unmarshal(*objmap["message"], &message)

                    logging.Error("taskRunner.webservices", "Authenticate", "(500) Unable to authenticate task runner", errors.New(message))
                }
            } else {
                logging.Debug("taskRunner.webservices", "Authenticate", "(" + strconv.Itoa(code) + ") Unable to authenticate task runner")
            }
        }
    }

    return TaskRunner{}
}
