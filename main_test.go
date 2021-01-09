package main

import (
	"fmt"
	"log"
	"os"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"))
	fmt.Println(os.Getenv("APP_DB_USERNAME"))
	fmt.Println(os.Getenv("APP_DB_PASSWORD"))
	fmt.Println(os.Getenv("APP_DB_NAME"))

	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM recipes")
	a.DB.Exec("ALTER SEQUENCE recipes_id_seq RESTART WITH 1")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS recipes
(
    id SERIAL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    CONSTRAINT recipes_pkey PRIMARY KEY (id)
)`
