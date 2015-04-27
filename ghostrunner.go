package main

import (
	"time"
	"ghostrunner/taskRunner"
	"ghostrunner/tasks"
    "ghostrunner/utils"
    "ghostrunner/logging"
    "github.com/kardianos/service"
)

type program struct {
	exit chan struct{}
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() error {
	ticker := time.NewTicker(3 * time.Minute)

	for {
		select {
		case _ = <-ticker.C:
			processing()
		case <-p.exit:
			ticker.Stop()
			return nil
		}
	}
	return nil
}

func (p *program) Stop(s service.Service) error {
	close(p.exit)
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "GhostRunner",
		DisplayName: "GhostRunner",
		Description: "GhostRunner task engine",
	}

	prg := &program{}
	
	s, err := service.New(prg, svcConfig)

	if err != nil {
        logging.Error("ghostrunner", "main", "Error creating service", err)
	}

	err = s.Run()

	if err != nil {
        logging.Error("ghostrunner", "main", "Error starting service", err)
	}
}

func processing() {	
	configuration, err := utils.LoadConfiguration()

	if err == nil {
		if (len(configuration.SessionId) == 0) {
			authenticate(&configuration)
		} else {
			returnCode := taskRunner.ValidateRunner(&configuration)

			switch returnCode {
			case 200: processTasks(&configuration)
			case 404: authenticate(&configuration)
			}
		}
	}
}

func authenticate(configuration *utils.Configuration) {
	sessionId := taskRunner.Authenticate(configuration)

	configuration.SessionId = sessionId

	utils.UpdateConfiguration(configuration)
	
	if (len(configuration.SessionId) > 0) {
		processTasks(configuration)
	}
}

func processTasks(configuration *utils.Configuration) {
	unprocessedTasks := taskRunner.GetTasks(configuration)

	tasks.ProcessTasks(configuration, unprocessedTasks)
}	