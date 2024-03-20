package messages

import "fmt"

func ErrorMessage(action string, err error) {
	msg := fmt.Sprintf("[!] Error to %s: %s", action, err)
	fmt.Println(msg)
}
