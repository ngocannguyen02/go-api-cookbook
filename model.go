package main

import (
	"database/sql"
)

type recipe struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (r *recipe) getRecipe(db *sql.DB) error {
	return db.QueryRow("SELECT title, description FROM recipes WHERE id=$1", r.ID).Scan(&r.Title, &r.Description)
}

func (r *recipe) updateRecipe(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE recipes SET title=$1, description=$2 WHERE id=$3",
			r.Title, r.Description, r.ID)

	return err
}

func (r *recipe) deleteRecipe(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM recipes WHERE id=$1", r.ID)

	return err
}

func (r *recipe) createRecipe(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO recipes(title, description) VALUES($1, $2) RETURNING id", r.Title, r.Description).Scan(&r.ID)

	if err != nil {
		return err
	}

	return nil
}

func getRecipes(db *sql.DB, start, count int) ([]recipe, error) {
	rows, err := db.Query(
		"SELECT id, title, description FROM recipes LIMIT $1 OFFSET $2", count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	recipes := []recipe{}

	for rows.Next() {
		var r recipe
		if err := rows.Scan(&r.ID, &r.Title, &r.Description); err != nil {
			return nil, err
		}
		recipes = append(recipes, r)
	}

	return recipes, nil
}
