package entities

import . "github.com/nonamecat19/go-orm/core/lib/entities"

type List struct {
	Model
	Name  string `db:"name" json:"name"`
	Items []Item `relation:"list_id" json:"items"`
}

func (list List) Info() string {
	return "lists"
}
