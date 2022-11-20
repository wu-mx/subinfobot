package handler

import (
	"fmt"
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", log.Lshortfile|log.Ldate|log.Ltime)

func HandleError(err error){
	if err != nil {
		logger.Println(fmt.Printf("Send message error: %s", err))
	}
}