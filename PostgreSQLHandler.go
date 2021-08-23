package main

import (
	"database/sql"
	"fmt"
)

const (
	host     = "appliancestatesdb.cyebc6nm0xm9.eu-west-2.rds.amazonaws.com"
	port     = 5432
	user     = "postgres"
	password = "asdbpassword"
	dbname   = "appliancestatesdb"
)

const (
	querySQLStatement        = `SELECT mode FROM HomeAppliances WHERE name = $1;`
	updateSingleSQLStatement = `UPDATE HomeAppliances SET mode = $2 WHERE name = $1;`
)

type PostgreSQLHandler struct {
	db* sql.DB
}

func (postgreHandler* PostgreSQLHandler) Connect() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s",
		host, port, user, password, dbname)

	var err error
	postgreHandler.db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
}

func (postgreHandler* PostgreSQLHandler) TestConnection(){
	result := postgreHandler.db.Ping()
	if result != nil {
		panic(result)
	}
}

func (postgreHandler* PostgreSQLHandler) CloseConnection(){
	err := postgreHandler.db.Close()
	if err != nil {
		panic(err)
	}
}

func (postgreHandler* PostgreSQLHandler) UpdateMode(applianceData map[string]interface{}) {
	_, err := postgreHandler.db.Exec(updateSingleSQLStatement, applianceData["Type"], applianceData["Mode"])
	if err != nil {
		panic(err)
	}
}

func QueryModeProp(db *sql.DB, appliance string) (mode string) {

	row := db.QueryRow(querySQLStatement, appliance)
	switch err := row.Scan(&mode); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		return mode
	default:
		panic(err)
	}
	return mode
}
