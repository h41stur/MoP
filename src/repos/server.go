// CLI server repositories
package repos

import (
	"MoP/src/messages"
	"MoP/src/models"
	"database/sql"
	"strconv"
)

type fieldAgents struct {
	db *sql.DB
}

func AgentRepo(db *sql.DB) *fieldAgents {
	return &fieldAgents{db}
}

func (repo fieldAgents) History(ID string, command []string) ([]models.PostCommandMessage, error) {
	agentID, _ := strconv.Atoi(ID)
	limit := 2

	if len(command) > 2 {
		limit, _ = strconv.Atoi(command[2])
	}

	lines, err := repo.db.Query(
		"select command, response from commands where response is not null and agent_id = ? limit ?",
		agentID,
		limit,
	)

	if err != nil {
		return nil, err
	}
	defer lines.Close()

	var commands []models.PostCommandMessage
	for lines.Next() {
		var command models.PostCommandMessage

		if err = lines.Scan(
			&command.Command,
			&command.Resp,
		); err != nil {
			return nil, err
		}
		commands = append(commands, command)
	}

	return commands, nil
}

func (repo fieldAgents) Search() ([]models.Agent, error) {
	lines, err := repo.db.Query(
		"select id, name, username, ifnull(alias, '') as alias, so, arch, hostname, ip, created from agents where active is true",
	)
	if err != nil {
		return nil, err
	}
	defer lines.Close()

	var agents []models.Agent
	for lines.Next() {
		var agent models.Agent

		if err = lines.Scan(
			&agent.ID,
			&agent.Name,
			&agent.Username,
			&agent.Alias,
			&agent.OS,
			&agent.Arch,
			&agent.Hostname,
			&agent.Ip,
			&agent.CreatedAt,
		); err != nil {
			return nil, err
		}

		agents = append(agents, agent)
	}

	return agents, nil
}

func (repo fieldAgents) SearchAgentTimeCreated() ([]models.Agent, error) {
	lines, err := repo.db.Query(
		"select id, created from agents where active is true",
	)
	if err != nil {
		return nil, err
	}
	defer lines.Close()

	var agents []models.Agent
	for lines.Next() {
		var agent models.Agent

		if err = lines.Scan(
			&agent.ID,
			&agent.CreatedAt,
		); err != nil {
			return nil, err
		}

		agents = append(agents, agent)
	}

	return agents, nil
}

func (repo fieldAgents) UpdateAlias(agentID int, alias string) error {
	statement, err := repo.db.Prepare(
		"update agents set alias = ? where id = ?",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(alias, agentID)
	if err != nil {
		return err
	}

	return nil
}

func (repo fieldAgents) KillAgent(ID int) error {
	statement, err := repo.db.Prepare(
		"update agents set active = false where id = ?",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(ID)
	if err != nil {
		return err
	}

	return nil

}

func (repo fieldAgents) SearchByID(ID uint64) (models.Agent, error) {
	line, err := repo.db.Query(
		"select id, name, username, ifnull(alias, '') as alias, so, arch, hostname, ip, created from agents where id = ? and active is true",
		ID,
	)
	if err != nil {
		return models.Agent{}, err
	}
	defer line.Close()

	var agent models.Agent

	if line.Next() {
		if err = line.Scan(
			&agent.ID,
			&agent.Name,
			&agent.Username,
			&agent.Alias,
			&agent.OS,
			&agent.Arch,
			&agent.Hostname,
			&agent.Ip,
			&agent.CreatedAt,
		); err != nil {
			return models.Agent{}, err
		}
	}

	return agent, nil
}

func (repo fieldAgents) CheckAgent(agentID uint64) (bool, error) {
	line, err := repo.db.Query(
		"select id from agents where id = ? and active is true",
		agentID,
	)
	if err != nil {
		return false, err
	}
	defer line.Close()

	for line.Next() {
		var agent models.Agent
		if err = line.Scan(
			&agent.ID,
		); err != nil {
			return false, err
		}
		if agent.ID > 0 {
			return true, nil
		}
	}

	return false, nil
}

func (repo fieldAgents) Command(ID uint64, command string) {
	var agent models.Agent
	agent, err := repo.SearchByID(ID)
	if err != nil {
		messages.ErrorMessage("search agent", err)
		return
	}

	statement, err := repo.db.Prepare(
		"insert into commands (agent_id, name, command) values (?, ?, ?)",
	)
	if err != nil {
		messages.ErrorMessage("input command in DB", err)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(agent.ID, agent.Name, command)
	if err != nil {
		messages.ErrorMessage("input command in DB", err)
		return
	}

}

func (repo fieldAgents) FileRequest(ID uint64, file string, fileName string, filePath string, action string) error {
	var agent models.Agent
	agent, err := repo.SearchByID(ID)
	if err != nil {
		return err
	}

	statement, err := repo.db.Prepare(
		"insert into files (agent_id, name, direction, filename, filepath, file) values (?, ?, ?, ?, ?, ?)",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(agent.ID, agent.Name, action, fileName, filePath, file)
	if err != nil {
		return err
	}

	return nil
}
