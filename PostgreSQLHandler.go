package main

import (
	"database/sql"
	"fmt"
)

type PostgreSQLHandler struct {
	db *sql.DB
}

const (
	host     = "appliancestatesdb.cyebc6nm0xm9.eu-west-2.rds.amazonaws.com"
	port     = 5432
	user     = "postgres"
	password = "asdbpassword"
	dbname   = "appliancestatesdb"
)

const updateSingleSQLStatement = `UPDATE HomeAppliances SET mode = $2 WHERE name = $1;`
const createToySQLStatement = `INSERT INTO HomeAppliances (name, mode) VALUES ($1, $2) ON CONFLICT DO NOTHING;`

func (postgreHandler *PostgreSQLHandler) Connect() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s",
		host, port, user, password, dbname)

	var err error
	postgreHandler.db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("db connection established")
}

func (postgreHandler *PostgreSQLHandler) TestConnection() {
	result := postgreHandler.db.Ping()
	if result != nil {
		panic(result)
	}
}

func (postgreHandler *PostgreSQLHandler) Disconnect() {
	err := postgreHandler.db.Close()
	if err != nil {
		panic(err)
	}
}

func (postgreHandler *PostgreSQLHandler) UpdateToyMode(toyName string, toyMode string) {
	_, err := postgreHandler.db.Exec(updateSingleSQLStatement, toyName, toyMode)
	if err != nil {
		fmt.Println("Couldnt update mode", err)
	}
}

func (postgreHandler *PostgreSQLHandler) CreateToy(toyName string, toyMode string) {
	_, err := postgreHandler.db.Exec(createToySQLStatement, toyName, toyMode)
	if err != nil {
		fmt.Println("unable to create toy object in db", err)
	}
}
