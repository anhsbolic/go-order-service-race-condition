package helper

import (
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

	f, err := os.OpenFile("logs/"+dt.Format("2006-01-02"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger := log.New(f, setLog.Status, log.LstdFlags)
	logger.Println(setLog.Data)
}
