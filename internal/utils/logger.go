package utils

import (
	"fmt"
)

// LogError is a simple function for error logging
func LogError(message string, err error) {
	if err != nil {
		fmt.Printf("%s: %v\n", message, err)
	} else {
		fmt.Println(message)
	}
}
