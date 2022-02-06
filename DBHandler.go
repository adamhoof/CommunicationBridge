package main

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
)

type DBHandler struct {
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

func (dbHandler *DBHandler) Connect() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s",
		host, port, user, password, dbname)

	var err error
	dbHandler.db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("db connection established")
}

func (dbHandler *DBHandler) TestConnection() {
	result := dbHandler.db.Ping()
	if result != nil {
		panic(result)
	}
}

func (dbHandler *DBHandler) Disconnect() {
	err := dbHandler.db.Close()
	if err != nil {
		panic(err)
	}
}

func (dbHandler *DBHandler) UpdateToyMode(toyName string, toyMode string) {
	_, err := dbHandler.db.Exec(updateSingleSQLStatement, toyName, toyMode)
	if err != nil {
		fmt.Println("Couldnt update mode", err)
	}
}

func (dbHandler *DBHandler) PullToyData() (toyBag map[string]Toy) {

	toyBag = make(map[string]Toy)

	rows, err := dbHandler.db.Query(toysDataQuery)
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
		toyAttributes := GeneralToy{}
		err = rows.Scan(&toyAttributes.name, pq.Array(&toyAttributes.availableModes), &toyAttributes.id, &toyAttributes.publishTopic, &toyAttributes.subscribeTopic)
		if err != nil {
			fmt.Println("unable to fetch toyAttributes data into toyAttributes", err)
		}
		toyBag[toyAttributes.Name()] = &toyAttributes
	}
	return toyBag
}
