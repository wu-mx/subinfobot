package handler

import (
	"fmt"
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", log.Lshortfile|log.Ldate|log.Ltime)

func HandleError(err error) {
	if err != nil {
		logger.Print(fmt.Sprintf("Send message error: %s", err))
	}
}
