// Agents controllers
package controllers

import (
	"MoP/src/db"
	"MoP/src/models"
	"MoP/src/repos"
	"MoP/src/responses"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// para cada enpoint da API uma função correspondente que utiliza o pacote responses para responder o solicitante

func CreateAgent(w http.ResponseWriter, r *http.Request) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Response(w, http.StatusUnprocessableEntity)
		return
	}

	// cria uma variável do tipo Agent pra ser preenchida a partir do pacote models
	var agent models.Agent
	if err = json.Unmarshal(requestBody, &agent); err != nil {
		responses.Response(w, http.StatusBadRequest)
		return
	}

	agent.Ip = r.RemoteAddr

	if err = agent.Prepare(); err != nil {
		responses.Response(w, http.StatusBadRequest)
		return
	}

	// cria a conexão com o banco a partir do pacote db
	db, err := db.Connect()
	if err != nil {
		responses.Response(w, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// cria o repositório a partir do pacote repo para executar os comando SQL
	repo := repos.NewAgentRepo(db)
	agent.ID, err = repo.Create(agent)
	if err != nil {
		responses.Response(w, http.StatusInternalServerError)
		return
	}

	responses.HandleNewAgent(agent)
	responses.Response(w, http.StatusCreated)
}

func PostFile(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	agentName := parameters["agentId"]
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Response(w, http.StatusUnprocessableEntity)
		return
	}

	// cria uma variável do tipo File pra ser preenchida a partir do pacote models
	var file models.File
	if err = json.Unmarshal(requestBody, &file); err != nil {
		responses.Response(w, http.StatusBadRequest)
		return
	}

	db, err := db.Connect()
	if err != nil {
		responses.Response(w, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// cria o repositório a partir do pacote repo para executar os comando SQL
	repo := repos.NewAgentRepo(db)
	err = repo.PostFile(file.ID, file.File)
	if err != nil {
		responses.Response(w, http.StatusInternalServerError)
		return
	}

	saveFile(file.Filename, file.File)
	repo.DeleteFile(agentName, file.ID)

	responses.Response(w, http.StatusCreated)
}

func saveFile(fileName string, file string) {
	d64, _ := base64.StdEncoding.DecodeString(file)
	err := os.WriteFile(fileName, d64, 0644)
	if err != nil {
		return
	}
}

func GetFile(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	agentName := parameters["agentId"]

	// cria a conexão com o banco a partir do pacote db
	db, err := db.Connect()
	if err != nil {
		responses.Response(w, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// cria o repositório a partir do pacote repo para executar os comando SQL
	repo := repos.NewAgentRepo(db)
	files, err := repo.GetFile(agentName)
	if err != nil {
		return
	}

	responses.JSON(w, http.StatusOK, files)
}

func AgentAlive(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	agentName := parameters["agentId"]

	db, err := db.Connect()
	if err != nil {
		responses.Response(w, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	repo := repos.NewAgentRepo(db)
	err = repo.UpdateAlive(agentName)
	if err != nil {
		return
	}

	responses.Response(w, http.StatusOK)

}

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	agentName := parameters["agentId"]
	fileID, _ := strconv.Atoi(parameters["fileId"])

	// cria a conexão com o banco a partir do pacote db
	db, err := db.Connect()
	if err != nil {
		responses.Response(w, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// cria o repositório a partir do pacote repo para executar os comando SQL
	repo := repos.NewAgentRepo(db)
	repo.DeleteFile(agentName, fileID)

	responses.Response(w, http.StatusOK)
}

func GetCommand(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	agentName := parameters["agentId"]

	// cria a conexão com o banco a partir do pacote db
	db, err := db.Connect()
	if err != nil {
		responses.Response(w, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// cria o repositório a partir do pacote repo para executar os comando SQL
	repo := repos.NewAgentRepo(db)
	commands, err := repo.GetCommand(agentName)
	if err != nil {
		return
	}

	responses.JSON(w, http.StatusOK, commands)
}

func PostCommand(w http.ResponseWriter, r *http.Request) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Response(w, http.StatusUnprocessableEntity)
		return
	}

	// cria uma variável do tipo PostCommandMessage pra ser preenchida a partir do pacote messages
	var command models.PostCommandMessage
	if err = json.Unmarshal(requestBody, &command); err != nil {
		responses.Response(w, http.StatusBadRequest)
		return
	}

	// cria a conexão com o banco a partir do pacote db
	db, err := db.Connect()
	if err != nil {
		responses.Response(w, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// cria o repositório a partir do pacote repo para executar os comando SQL
	repo := repos.NewAgentRepo(db)
	err = repo.PostCommand(command.ID, command.Resp)
	if err != nil {
		responses.Response(w, http.StatusInternalServerError)
		return
	}

	responses.HandleCommand(command)
	responses.Response(w, http.StatusCreated)
}
