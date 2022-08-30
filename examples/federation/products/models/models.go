package models

type Product struct {
	Upc     string `json:"upc"`
	Name    string `json:"name"`
	Price   int    `json:"price"`
	InStock int    `json:"inStock"`
}

func (Product) IsEntity() {}

var Hats = []*Product{
	{
		Upc:     "top-1",
		Name:    "Trilby",
		Price:   11,
		InStock: 500,
	},
	{
		Upc:     "top-2",
		Name:    "Fedora",
		Price:   22,
		InStock: 1200,
	},
	{
		Upc:     "top-3",
		Name:    "Boater",
		Price:   33,
		InStock: 850,
	},
}
