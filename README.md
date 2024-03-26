# Master of Puppets - MoP

![Master of Puppets](https://raw.githubusercontent.com/h41stur/MoP/main/server/resources/master.png)

Master of Puppets is a personal project to build a C2 (Command and Control) focused on asynchronous communication via HTTP requests.

The main objective was to learn in the process as well as to establish knowledge, both in the functioning of a C2 and in the Golang programming language.

This project will still undergo several updates and implementations along the way.

## Running the Project

1. Before building the server and agent, make the appropriate configurations in the Config package (config/config.go), such as:

- Server
- Hostname

2. It is also necessary to start the MySQL database Docker by executing the sql/start-db.sh script.

```bash
cd sql
./create-db.sh
```

3. Generate the certificates from TLS communication.

```bash
cd server/resources
./generate-cert.sh
```

4. Compile the server.

```bash
cd server
go build MoP.go
./Mop
```


## Building Agents

To build an agent to Windows machines on server terminal:

```bash
GOOS=windows go build -ldflags -H=windowsgui agent.go
```

To build an agent to Linux machines on server terminal:

```bash
GOOS=linux go build agent.go
```

To build agents to Linux and Windows machines on MoP terminal:

```bash
MoP> build <name>
```

This command will auto build the agents and put them ina a file server at https://server/downloads/  (default name "agent").

## Commands

```
	SERVER COMMANDS:

		!<command>: run SO command on server machine.

		alias <alias>: set an alias to an agent when one that agent is selected.

		build [name]: build Linux and Windows agents and put them on a File Server at https://server/drop/agents/ (default name "agent").

		help: show this help.
	
		show agents: list the active agents to interact.

		show modules: list the modules to inject on the agent machine.
			
		select <id>: select an agent by id to interact.


	AGENT COMMANDS:

		download <filepath>: download a file from agent machine.

		persist: make a copy of the agent to startup directory on Windows machines.

		ps: shows the process list running on agent machine.

		screenshot: take a screenshot from agent machine.

		shell [port]: initiate a TCP reverse shell listening on informated port (default 8081). (very easy to detect and not recomended)

		show history: list the command history when one agent is selected.

		sleep <time seconds>: change the sleep time before send response to server (default 10).

		upload <file path>: send a file to agent machine.

		use <module name>: inject the scripts listed on "show modules" command in the agent machine.
```

