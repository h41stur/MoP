package messages

import (
	"fmt"

	"github.com/fatih/color"
)

// colors
var (
	Red        = color.New(color.FgRed)
	RedBold    = Red.Add(color.Bold)
	Blue       = color.New(color.FgBlue)
	BlueBold   = Blue.Add(color.Bold)
	White      = color.New(color.FgWhite)
	WhiteBold  = White.Add(color.Bold)
	Green      = color.New(color.FgGreen)
	GreenBold  = Green.Add(color.Bold)
	Yellow     = color.New(color.FgYellow)
	YellowBold = Yellow.Add(color.Bold)
)

func Banner(hostname string) {
	Red.Println("+===============================================================================+")
	Red.Println("+                                MASTER OF PUPPETS                              +")
	Red.Println("+===============================================================================+")
	WhiteBold.Printf("================ Runing on: %s\n\n", hostname)
}

func ErrorMessage(action string, err error) {
	RedBold.Printf("\n\n[!] Error to %s: %s\n", action, err)
	fmt.Println()
}

func SelectAgentMessage() {
	RedBold.Println("\nYou need to select an agent first!")
	fmt.Println()
}

func WarningMessage(action string) {
	YellowBold.Printf("\n\n[!] %s\n\n", action)
}

func OkMessage(action string) {
	GreenBold.Printf("\n\n[!] %s\n\n", action)
}