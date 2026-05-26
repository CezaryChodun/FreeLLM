package models

type Model struct {
	ID       int    `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Provider string `db:"provider" json:"provider"`
	Instance int    `db:"instance" json:"instance"`
}
