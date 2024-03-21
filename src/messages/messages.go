package messages

import (
	"fmt"

	"github.com/fatih/color"
)

// colors
var (
	Red       = color.New(color.FgRed)
	RedBold   = Red.Add(color.Bold)
	Blue      = color.New(color.FgBlue)
	BlueBold  = Blue.Add(color.Bold)
	White     = color.New(color.FgWhite)
	WhiteBold = White.Add(color.Bold)
	Green     = color.New(color.FgGreen)
	GreenBold = Green.Add(color.Bold)
)

func Banner(hostname string) {
	Red.Println("+===============================================================================+")
	Red.Println("+                                MASTER OF PUPPETS                              +")
	Red.Println("+===============================================================================+")
	WhiteBold.Printf("================ Runing on: %s\n\n", hostname)
}

func ErrorMessage(action string, err error) {
	msg := fmt.Sprintf("[!] Error to %s: %s", action, err)
	fmt.Println(msg)
}
