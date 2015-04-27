package tasks

import (
	"time"
)

type Task struct {
    ExternalId string `json:"externalId"`
    EncryptedScripts []string `json:"scripts"`
    Scripts []TaskScript
}

type TaskScript struct {
	Id int `json:"id"`
	Type string `json:"type"`
	Content string `json:"content"`
	Log string
	Position int `json:"position"`
	Status string
	Started time.Time
	Completed time.Time
}