package main

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
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

const updateSingleSQLStatement = `UPDATE HomeAppliances SET current_mode = $2 WHERE name = $1;`
const toysDataQuery = `SELECT name, available_modes, id, publish_topic, subscribe_topic FROM HomeAppliances;`

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

func (postgreHandler *PostgreSQLHandler) PullToyData(toyBag map[string]Toy) {

	rows, err := postgreHandler.db.Query(toysDataQuery)
	if err != nil {
		fmt.Println("unable to query data", err)
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			fmt.Println("unable to close rows", err)
		}
	}(rows)

	for rows.Next() {
		toy := ToyAttributes{}
		err = rows.Scan(&toy.name, pq.Array(&toy.availableModes), &toy.id, &toy.publishTopic, &toy.subscribeTopic)
		if err != nil {
			fmt.Println("unable to fetch toy data into toyAttributes", err)
		}
		toyBag[toy.Name()] = &toy
	}
}
