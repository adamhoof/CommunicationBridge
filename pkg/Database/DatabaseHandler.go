package database

import connectable "CommunicationBridge/pkg/ConnectableDevices"

type DatabaseHandler interface {
	Connect(config string) error
	Ping() error
	RegisterToy(toy *connectable.Toy) (err error)
	PullToyData(toys map[string]*connectable.Toy)
}
