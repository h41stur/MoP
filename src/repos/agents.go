// Agents repositories
package repos

import (
	"MoP/src/models"
	"database/sql"
	"fmt"
)

type agents struct {
	db *sql.DB
}

func NewAgentRepo(db *sql.DB) *agents {
	return &agents{db}
}

func (repo agents) Create(agent models.Agent) (uint64, error) {
	// checa se o agente já é cadastrado e está desativado
	status, id := repo.checkKilled(agent.Name)

	// se estiver desativado, reativa
	if status {
		statement, err := repo.db.Prepare(
			"update agents set active = true, created = current_timestamp() where name = ?",
		)
		if err != nil {
			return 0, err
		}
		defer statement.Close()

		_, err = statement.Exec(agent.Name)
		if err != nil {
			return 0, err
		}

		return uint64(id), nil

	} else { // se for novo, cadastra
		statement, err := repo.db.Prepare(
			"insert into agents (name, username, so, arch, hostname, ip) values (?, ?, ?, ?, ?, ?)",
		)
		if err != nil {
			return 0, err
		}
		defer statement.Close()

		result, err := statement.Exec(agent.Name, agent.Username, agent.SO, agent.Arch, agent.Hostname, agent.Ip)
		if err != nil {
			return 0, err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return 0, err
		}

		return uint64(lastID), nil
	}
}

func (repo agents) checkKilled(name string) (bool, int) {
	line, err := repo.db.Query(
		"select id, active from agents where name = ?",
		name,
	)
	if err != nil {
		return false, 0
	}
	defer line.Close()

	var kill struct {
		id     int
		status bool
	}
	for line.Next() {
		if err = line.Scan(
			&kill.id,
			&kill.status,
		); err != nil {
			return false, 0
		}
	}

	if kill.id == 0 {
		return false, 0
	}

	if bool(kill.status) {
		return false, 0
	}

	return true, kill.id

}

func (repo agents) GetFile(agentName string) ([]models.File, error) {
	lines, err := repo.db.Query(
		"select id, filename, filepath, file, direction from files where name = ?",
		agentName,
	)
	if err != nil {
		return nil, err
	}
	defer lines.Close()

	var files []models.File
	for lines.Next() {
		var file models.File
		if err = lines.Scan(
			&file.ID,
			&file.Filename,
			&file.FilePath,
			&file.File,
			&file.Direction,
		); err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	return files, nil
}

func (repo agents) GetCommand(agentName string) ([]models.Command, error) {
	lines, err := repo.db.Query(
		"select id, command from commands where name = ? and response is null",
		agentName,
	)
	if err != nil {
		return nil, err
	}
	defer lines.Close()

	var commands []models.Command
	for lines.Next() {
		var command models.Command
		if err = lines.Scan(
			&command.ID,
			&command.Command,
		); err != nil {
			return nil, err
		}

		commands = append(commands, command)
	}

	return commands, nil
}

func (repo agents) UpdateAlive(agentName string) error {
	statement, err := repo.db.Prepare(
		"update agents set created = current_timestamp(), active = 1 where name = ?",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(agentName)
	if err != nil {
		return err
	}

	return nil
}

func (repo agents) PostCommand(ID int, resp string) (int, error) {
	statement, err := repo.db.Prepare(
		"update commands set response = ? where id = ?",
	)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	_, err = statement.Exec(resp, ID)
	if err != nil {
		return 0, err
	}

	lines, err := repo.db.Query(
		"select agent_id from commands where id = ?",
		ID,
	)
	if err != nil {
		return 0, err
	}
	defer lines.Close()

	var agentID int
	for lines.Next() {
		
		if err = lines.Scan(
			&agentID,
		); err != nil {
			fmt.Println(err)
			return 0, err
		}
	}

	return agentID, nil

}

func (repo agents) PostFile(ID int, file string) error {
	statement, err := repo.db.Prepare(
		"update files set file = ? where id = ?",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(file, ID)
	if err != nil {
		return err
	}

	return nil

}

func (repo agents) DeleteFile(agentName string, fileID int) {
	statement, err := repo.db.Prepare(
		"delete from files where id = ? and name = ?",
	)
	if err != nil {
		return
	}
	defer statement.Close()

	_, err = statement.Exec(fileID, agentName)
	if err != nil {
		return
	}
}
