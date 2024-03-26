CREATE DATABASE IF NOT EXISTS MoP;
USE MoP;

DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS commands;
DROP TABLE IF EXISTS agents;

CREATE TABLE agents(
    id int auto_increment primary key,
    name varchar(50) not null unique,
    username varchar(50) not null,
    alias varchar(50) unique,
    active boolean not null default 1,
    so varchar(20) not null,
    arch varchar(20) not null,
    hostname varchar(50) not null,
    ip varchar(20) not null,
    created timestamp not null default current_timestamp()
) ENGINE=INNODB;

CREATE TABLE commands(
    id int auto_increment primary key,
    agent_id int,
    name varchar(50) not null,
    command varchar(255) not null,
    response longtext,
    created timestamp not null default current_timestamp(),
    FOREIGN KEY (agent_id) REFERENCES agents(id) ON DELETE CASCADE
) ENGINE=INNODB;

CREATE TABLE files(
    id int auto_increment primary key,
    agent_id int,
    name varchar(50) not null,
    direction varchar(10) not null,
    filename varchar(50) not null,
    filepath varchar(80) not null,
    file longtext,
    created timestamp not null default current_timestamp(),
    FOREIGN KEY (agent_id) REFERENCES agents(id) ON DELETE CASCADE
) ENGINE=INNODB;
