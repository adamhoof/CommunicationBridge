package database

import (
	"RPICommandHandler/pkg/ConnectableDevices"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
)

type PostgresHandler struct {
	db *sql.DB
}

const toysDataQuery = `SELECT name, available_modes, id, publish_topic, subscribe_topic FROM HomeAppliances;`

func (handler *PostgresHandler) Connect(connectionString *string) (err error) {
	handler.db, err = sql.Open("postgres", *connectionString)
	if err != nil {
		return err
	}
	return handler.Ping()
}

func (handler *PostgresHandler) Ping() (err error) {
	return handler.db.Ping()
}

func (handler *PostgresHandler) Disconnect() {
	err := handler.db.Close()
	if err != nil {
		panic(err)
	}
}

func (handler *PostgresHandler) PullToyData(toyBag map[string]*connectable.Toy) {

	rows, err := handler.db.Query(toysDataQuery)
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
		toy := connectable.Toy{}
		err = rows.Scan(&toy.Name, pq.Array(&toy.AvailableModes), &toy.Id, &toy.PublishTopic, &toy.SubscribeTopic)
		if err != nil {
			fmt.Println("unable to fetch toy data", err)
		}
		toyBag[toy.Name] = &toy
	}
}
