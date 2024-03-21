// CLI server responses
package responses

import (
	"MoP/src/messages"
	"MoP/src/models"
	"encoding/base64"
	"fmt"
)

func PrintAgents(data []models.Agent) {

	for _, agent := range data {
		fmt.Println("\nAgent ID: ", agent.ID)
		fmt.Println("Agent name: ", agent.Name)
		fmt.Println("Agent username: ", agent.Username)
		fmt.Println("Agent alias: ", agent.Alias)
		fmt.Println("Agent SO: ", agent.SO)
		fmt.Println("Agent arch: ", agent.Arch)
		fmt.Println("Agent hostname: ", agent.Hostname)
		fmt.Println("Agent address: ", agent.Ip)
		fmt.Println("Agent created at: ", agent.CreatedAt)
		fmt.Println()
	}
}

func PrintAgent(agent models.Agent) {

	fmt.Println("\nAgent ID: ", agent.ID)
	fmt.Println("Agent name: ", agent.Name)
	fmt.Println("Agent username: ", agent.Username)
	fmt.Println("Agent alias: ", agent.Alias)
	fmt.Println("Agent SO: ", agent.SO)
	fmt.Println("Agent arch: ", agent.Arch)
	fmt.Println("Agent hostname: ", agent.Hostname)
	fmt.Println("Agent address: ", agent.Ip)
	fmt.Println()
}

func HandleNewAgent(agent models.Agent) {
	messages.GreenBold.Println("\n\nNew agent connected:")
	PrintAgent(agent)
}

func HandleCommand(resp models.PostCommandMessage) {
	response, _ := base64.StdEncoding.DecodeString(resp.Resp)
	fmt.Println()
	fmt.Println("Command: ", resp.Command)
	fmt.Println("Response: \n\n", string(response))
	fmt.Println()
}
