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

const toysDataQuery = `SELECT name, available_modes, room, id, publish_topic, subscribe_topic FROM HomeAppliances where room = $1;`
const placesDataQuery = `SELECT * FROM places;`

func (handler *PostgresHandler) Connect(connectionString string) (err error) {
	handler.db, err = sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}
	return handler.Ping()
}

func (handler *PostgresHandler) Ping() (err error) {
	return handler.db.Ping()
}

func (handler *PostgresHandler) ExecuteStatement(statement string) {
	if _, err := handler.db.Exec(statement); err != nil {
		fmt.Println(err)
	}
}

func (handler *PostgresHandler) Disconnect() {
	if err := handler.db.Close(); err != nil {
		fmt.Println("failed to close db connection", err)
	}
}

func (handler *PostgresHandler) PullAvailableRooms() (places []string) {
	rows, err := handler.db.Query(placesDataQuery)
	if err != nil {
		fmt.Println("unable to query places data", err)
	}
	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			fmt.Println("unable to close rows", err)
		}
	}(rows)

	for rows.Next() {
		var place string
		if err = rows.Scan(&place); err != nil {
			fmt.Println("unable to fetch places data", err)
		}
		places = append(places, place)
	}
	return places
}

func (handler *PostgresHandler) PullToyDataBasedOnRoom(toyBag map[string]*connectable.Toy, room string) {
	rows, err := handler.db.Query(toysDataQuery, room)
	if err != nil {
		fmt.Println("unable to query data", err)
	}
	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			fmt.Println("unable to close rows", err)
			return
		}
	}(rows)

	for rows.Next() {
		toy := connectable.Toy{}
		if err = rows.Scan(&toy.Name, pq.Array(&toy.AvailableModes), &toy.Room, &toy.Id, &toy.PublishTopic, &toy.SubscribeTopic); err != nil {
			fmt.Println("unable to fetch toy data", err)
		}
		toyBag[toy.Name] = &toy
	}
}
