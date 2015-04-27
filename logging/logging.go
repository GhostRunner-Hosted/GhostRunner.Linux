package logging

import (
	"fmt"
	"os"
    "log"
    "io/ioutil"
)

var config Configuration 

func init() {
	config = LoadConfiguration()
}

func Debug(script string, method string, message string) {
	if config.DebugEnabled {
		f,err := os.OpenFile(config.DebugLogLocation, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		
		if err != nil {
		    fmt.Println("error opening file '" + config.DebugLogLocation + "': ", err)
		} else {
			defer f.Close()
			defer checkFileSize(config.DebugLogLocation, config.DebugMaxFileSizeKB)

			log.SetOutput(f)

			log.Println(script + ": " + method + ": ", message)
		}
	}
}

func Error(script string, method string, message string, err error) {
	Debug(script, method, message)

			    fmt.Println(config.ErrorLogLocation)

	f,fileErr := os.OpenFile(config.ErrorLogLocation, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	
	if fileErr != nil {
	    fmt.Println("error opening file '" + config.ErrorLogLocation + "': ", fileErr)
	} else {
		defer f.Close()
		defer checkFileSize(config.ErrorLogLocation, config.ErrorMaxFileSizeKB)

		log.SetOutput(f)

		log.Println(script + ": " + method + ": " + ": " + message, err)
	}
}

func checkFileSize(fileName string, maxFileSize int) {
    data, err := ioutil.ReadFile(fileName)

	if err != nil {
	    fmt.Println("Error opening file to check size  at '" + fileName + "': ", err)
	    return
	}

	bytesAllowed := maxFileSize * 1024

	if (len(data) > 0) && (len(data) > bytesAllowed) {
		os.Remove(fileName)

		outputData := data[(len(data)-bytesAllowed):(len(data)-1)]

	    loggingFile, err := os.Create(fileName)

	    if err != nil {
		    fmt.Println("Error creating node script at " + fileName)
	        return 
	    }
	    
	    defer loggingFile.Close()

	    _, err = loggingFile.Write(outputData)
	}
}