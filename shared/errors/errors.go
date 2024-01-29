package errors

import (
	"fmt"
	"log"
)

func HandleError(err error) {
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error: %s", err))
	}
}
