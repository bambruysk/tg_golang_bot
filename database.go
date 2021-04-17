package main

import (
	_ "database/sql"

	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// // Holde - basic structures  for holde
// type Holde struct {
// 	// Имя поместья
// 	Name string
// 	// ID -  номер поместья
// 	ID int
// 	// Money - текущее накопелния
// 	Amount Money

// 	// Уровень
// 	Level int

// 	// Владелец
// 	Owner string

// 	//
// 	LastVisit time.Time
// }

var schema = `
CREATE TABLE IF NOT EXISTS player (
    name text
);

CREATE TABLE IF NOT EXISTS user (
    country text,
    city text NULL,
    telcode integer
)

CREATE TABLE IF NOT EXISTS holde (
	name text,
	id int not null,
	level int null,
	Amount int not null,
	LastVisit timestamp,
	FOREIGN KEY (player) owner NULL, 
	UNIQUE(name, id)
)

`

func ConnectDB () error {
	    // this Pings the database trying to connect
    // use sqlx.Open() for sql.Open() semantics
    db, err := sqlx.Connect("sqlite3", "__deleteme.db")
    if err != nil {
		log.Fatalln(err)
		return err
    }

    // exec the schema or fail; multi-statement Exec behavior varies between
    // database drivers;  pq will exec them all, sqlite3 won't, ymmv
	db.MustExec(schema)
	
	tx := db.MustBegin()

    // Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person
    res, err := tx.NamedExec("INSERT INTO player (name) VALUES (:name )", &Player{
    	Name:    "fedor",
	})

	if err !=  nil {
		log.Fatalln(err)
		return err
	}

	log.Println(res)
	
	tx.Commit()
	
	return nil

} 