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
	SO        string    `json:"so,omitempty"`
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

func (agent *Agent) Prepare() error {
	if err := agent.validate(); err != nil {
		return err
	}

	agent.format()
	return nil
}

func (agent *Agent) validate() error {
	if agent.Name == "" {
		return errors.New("Empty name")
	}

	if agent.Username == "" {
		return errors.New("Empty username")
	}

	if agent.SO == "" {
		return errors.New("Empty SO")
	}

	if agent.Hostname == "" {
		return errors.New("Empty hostname")
	}

	return nil
}

func (agent *Agent) format() {
	agent.Name = strings.TrimSpace(agent.Name)
	agent.Username = strings.TrimSpace(agent.Username)
	agent.SO = strings.TrimSpace(agent.SO)
	agent.Hostname = strings.TrimSpace(agent.Hostname)
	agent.Arch = strings.TrimSpace(agent.Arch)
}
