package database

import "RPICommandHandler/pkg/ConnectableDevices"

type DatabaseHandler interface {
	Connect(config string) error
	Ping() error
	ExecuteStatement(statement string)
	PullToyData(toys map[string]*connectable.Toy)
}
