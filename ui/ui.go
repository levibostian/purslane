package ui

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// HandleError pass in error and we will handle it.
func HandleError(err error) {
	if err != nil {
		Error("\nError encountered!")
		fmt.Println(err)
		os.Exit(1)
	}
}

func Abort(message string) {
	Error(message)
	os.Exit(1)
}

// Error show a message in red
func Error(message string) {
	color.Red(message)
}

// Message Show a neutral message in white
func Message(message string) {
	color.White(message)
}
