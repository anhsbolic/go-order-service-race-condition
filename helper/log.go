package helper

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Logging struct {
	Status string      `json:"string"`
	Data   interface{} `json:"data"`
}

func Log(status string, data interface{}) {
	dt := time.Now()

	setLog := Logging{
		Status: status,
		Data:   data,
	}

	path := fmt.Sprintf(`/logs/%s`, dt.Format("2006-01-02"))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModeDir)
	}

	f, _ := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModeDir)
	defer f.Close()

	logger := log.New(f, setLog.Status, log.LstdFlags)
	logger.Println(setLog.Data)
}
