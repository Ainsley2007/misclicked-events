package utils

import (
	"fmt"
)

func LogError(message string, err error) {
	if err != nil {
		fmt.Printf("%s: %v\n", message, err)
	} else {
		fmt.Println(message)
	}
}
