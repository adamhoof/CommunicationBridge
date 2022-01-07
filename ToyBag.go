package main

type ToyStorage struct {
	bag map[string]Toy
}

func (toyStorage *ToyStorage) fill(toys []Toy) {
	toyStorage.bag = make(map[string]Toy)

	for _, toy := range toys {
		toyStorage.bag[toy.Name()] = toy
		toy.Name()
	}
}
