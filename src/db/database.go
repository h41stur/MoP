package db

import (
	"MoP/src/config"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() (*sql.DB, error) {
	config := config.Load()
	db, err := sql.Open("mysql", config.DbConnectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func CheckDB() {
	db, err := Connect()
	if err != nil {
		fmt.Println(`
Database is down!
If "sudo docker ps -a" command don't show "MoPDB" container exited
Run the command:
"sql/create-db.sh"
If "sudo docker ps -a" command show "MoPDB" container exited
Run the command:
sudo docker start MoPDB`)
		os.Exit(0)
	}
	defer db.Close()

	tables := [3]string{"agents", "commands", "files"}

	for _, table := range tables {
		query := "select 1 from " + table
		_, err = db.Query(query)
		if err != nil {
			fmt.Printf("Table %s is not there!\n", table)
			fmt.Println(`
Run the command:
sudo docker exec -it MoPDB mysql -e 'source /scripts/sql.sql'`)
			os.Exit(0)
		}
	}

}
