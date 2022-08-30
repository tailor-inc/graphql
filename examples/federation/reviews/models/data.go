package models

var Reviews = []*Review{
	{
		Body:    "A highly effective form of birth control.",
		Product: &Product{Upc: "top-1"},
		Author:  &User{ID: "1234"},
	},
	{
		Body:    "Fedoras are one of the most fashionable hats around and can look great with a variety of outfits.",
		Product: &Product{Upc: "top-2"},
		Author:  &User{ID: "1234"},
	},
	{
		Body:    "This is the last straw. Hat you will wear. 11/10",
		Product: &Product{Upc: "top-3"},
		Author:  &User{ID: "7777"},
	},
}
