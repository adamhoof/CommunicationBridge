package main

type ToyBag struct {
	bag map[string]Toy
}

func (toyBag *ToyBag) fill(toys []Toy) {
	toyBag.bag = make(map[string]Toy)

	for _, toy := range toys {
		toyBag.bag[toy.Name()] = toy
		toy.Name()
	}
}
