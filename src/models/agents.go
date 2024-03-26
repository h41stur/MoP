package models

import (
	"errors"
	"strings"
	"time"
)

type Agent struct {
	ID        uint64    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Username  string    `json:"username,omitempty"`
	Alias     string    `json:"alias,omitempty"`
	OS        string    `json:"os,omitempty"`
	Arch      string    `json:"arch,omitempty"`
	Hostname  string    `json:"hostname,omitempty"`
	Ip        string    `json:"ip,omitempty"`
	CreatedAt time.Time `json:"created,omitempty"`
}

type Command struct {
	ID      int    `json:"id"`
	Command string `json:"command"`
}

type File struct {
	ID        int    `json:"id"`
	Filename  string `json:"filename"`
	FilePath  string `json:"filepath"`
	File      string `json:"file,omitempty"`
	Direction string `json:"direction"`
}

type NewAgent struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Hostname string `json:"hostname"`
}

type PostCommandMessage struct {
	ID      int    `json:"id"`
	Command string `json:"command"`
	Resp    string `json:"resp"`
}

var Modules = map[string]map[string]map[string]string{
	"Windows": {
		"Local Privilege Escalation": {
			"PowerUp": "PowerUp aims to be a clearinghouse of common Windows privilege escalation vectors that rely on misconfigurations.",
		},
		"Domain Enumeration": {
			"PowerView": "PowerView is a PowerShell tool to gain network situational awareness on Windows domains. It contains a set of pure-PowerShell replacements for various windows \"net *\" commands, which utilize PowerShell AD hooks and underlying Win32 API functions to perform useful Windows domain functionality.",
			"PowerView_dev": "Alternative to PowerView",
			"PowerView_2022": "Alternative to PowerView with support to Windows Server 2022",
			"SharpHound": "A BloodHound PowerShell collector",
		},
	},
	"Linux": {
		"": {
			"": "",
		},
	},
}


func (agent *Agent) Prepare() error {
	if err := agent.validate(); err != nil {
		return err
	}

	agent.format()
	return nil
}

func (agent *Agent) validate() error {
	if agent.Name == "" {
		return errors.New("empty name")
	}

	if agent.Username == "" {
		return errors.New("empty username")
	}

	if agent.OS == "" {
		return errors.New("empty SO")
	}

	if agent.Hostname == "" {
		return errors.New("empty hostname")
	}

	return nil
}

func (agent *Agent) format() {
	agent.Name = strings.TrimSpace(agent.Name)
	agent.Username = strings.TrimSpace(agent.Username)
	agent.OS = strings.TrimSpace(agent.OS)
	agent.Hostname = strings.TrimSpace(agent.Hostname)
	agent.Arch = strings.TrimSpace(agent.Arch)
}
