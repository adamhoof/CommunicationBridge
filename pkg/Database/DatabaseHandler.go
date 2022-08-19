package database

import connectable "CommunicationBridge/pkg/ConnectableDevices"

type DatabaseHandler interface {
	Connect(config string) error
	Ping() error
	RegisterToy(toy *connectable.Toy) (err error)
	UpdateDeviceIpAddress(ipToSet string, name string) (err error)
	PullToyData(toys map[string]*connectable.Toy)
}
