# Master of Puppets - MoP

Master of Puppets is a personal project to build a C2 (Command and Control) focused on asynchronous communication via HTTP requests.

The main objective was to learn in the process as well as to establish knowledge, both in the functioning of a C2 and in the Golang programming language.

This project will still undergo several updates and implementations along the way.

## Running the Project

Before building the server and agent, make the appropriate configurations in the Config package (config/config.go), such as:

- Server
- Hostname

It is also necessary to start the MySQL database Docker by executing the sql/start-db.sh script.

```bash
./create-db.sh
```

## Building Agents

To build an agent to Windows machines:

```bash
GOOS=windows go build -ldflags -H=windowsgui agent.go
```

To build an agent to Linux machines:

```bash
GOOS=linux go build agent.go
```

## Commands

```
	SERVER COMMANDS:

		!<command>: run SO command on server machine.

		alias <alias>: set an alias to an agent when one that agent is selected.

		build <name>: build Linux and Windows agents and put them on a File Server at http(s)://server/downloads/ (default name "agent").

		help: show this help.
	
		show agents: list the active agents to interact.
			
		select <id>: select an agent by id to interact.


	AGENT COMMANDS:

		download <filepath>: download a file from agent machine.

		persist: make a copy of the agent to startup directory on Windows machines.

		ps: shows the process list running on agent machine.

		screenshot: take a screenshot from agent machine.

		show history: list the command history when one agent is selected.

		sleep <time seconds>: change the sleep time before send response to server (default 10).

		upload <file path>: send a file to agent machine.
```

