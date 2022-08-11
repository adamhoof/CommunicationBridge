package database

import "RPICommandHandler/pkg/ConnectableDevices"

type DatabaseHandler interface {
	Connect(config string) error
	Ping() error
	ExecuteStatement(statement string)
	PullToyDataBasedOnRoom(toys map[string]*connectable.Toy, room string)
}
