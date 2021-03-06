package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
)

//
type Location struct {
	Name string
}

var schema = `

CREATE TABLE IF NOT EXISTS location (id serial PRIMARY KEY, name varchar unique);

CREATE TABLE IF NOT EXISTS players (
    id serial PRIMARY KEY,
    name varchar,
    amount decimal(10, 2),
    last_visit TIMESTAMP
);


CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
	chat_id numeric unique,
    name varchar(80) not null,

    location int,

    FOREIGN KEY (location) REFERENCES location (id) ON DELETE
    SET
        NULL
);

CREATE TABLE IF NOT EXISTS holdes (
    id serial PRIMARY KEY,
    name varchar NOT NULL,
    amount decimal,
    level integer not null default 1,
    owner_id integer REFERENCES players (id) ON DELETE SET NULL ,
    last_visit timestamp 
);


`

// func ConnectDB() *sqlx.DB {
// 	db, err := sqlx.Connect("pgx", "user=holde_tg_bot passwd=holde_tg_bot  dbname=test sslmode=disable")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	return db
// }

func ConnectDB() (*pgx.Conn, error) {

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
	}
	return conn, err
}

// Users DB

//GetFromDB read user from db
func (u *User) ReadFromDB(conn *pgx.Conn) error {

	row := conn.QueryRow(context.Background(), "select name from users where chat_id=$1", u.ChatID)
	fmt.Println(" Row found is", row)
	row.Scan(u.Name)
	return nil
}

func (u *User) CreateDB(conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(),
		`insert into users(name, chat_id, location) 
		values($1, $2, (select id from location where name=$3))  
		on conflict (chat_id) 
		do update set name=$1 ,location=(select id from location where name=$3) `,
		u.Name, u.ChatID, u.Location)
	return err
}

func (u *User) UpdateDB(conn *pgx.Conn) error {

	_, err := conn.Exec(context.Background(),
		"update set name=$1 ,location=(select id from location where name=$3) where chat_id=$2",
		u.Name, u.ChatID, u.Location)
	if err != nil {
		log.Fatalf(" update to user faild %v", err)
	}

	return nil
}

func (u *User) ReadDB(conn *pgx.Conn) error {

	rows, err := conn.Query(context.Background(),
		"select users.name, chat_id, location.name from users join location on users.location = location.id where chat_id=$1", u.ChatID)
	if err != nil {
		return err
	}
	if !rows.Next() {
		return errors.New("User not found")
	}
	rows.Scan(u.Name, u.ChatID, u.Location)

	return nil
}

func (l *Location) AddDB(conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(), "insert into location(name) values($1) on conflict do nothing", l.Name)
	return err
}

func LocationsReadFromDB(conn *pgx.Conn) ([]Location, error) {
	rows, err := conn.Query(context.Background(), "select name  from location ")

	locs := []Location{}

	for rows.Next() {
		var loc Location
		rows.Scan(&loc.Name)
		locs = append(locs, loc)
	}

	return locs, err
}

func (l *Location) DeleteDB(conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(), "delete from location where name=$1", l.Name)
	return err
}

// In memory
type UsersDB struct {
	conn *pgx.Conn
}

func NewUsersDB() (UsersDB, error) {
	conn, err := ConnectDB()
	if err != nil {
		log.Fatal("Database connect error")
	}
	return UsersDB{
		conn: conn,
	}, nil
}

func (u UsersDB) Get(id UserID) (User, error) {

	// Get user by chat_id
	user := User{
		ChatID: int(id),
	}
	if err := user.ReadDB(u.conn); err != nil {
		return user, err
	}
	return user, nil
}

func (u UsersDB) GetOrCreate(id UserID) (User, bool) {

	user := User{
		ChatID: int(id),
	}

	rows, err := u.conn.Query(context.Background(),
		"select users.name, chat_id, location.name from users join location on users.location = location.id where chat_id=$1", user.ChatID)
	if err != nil {
		log.Fatalf("Database read error")
	}
	if !rows.Next() {
		log.Print("Create new user")
		if err := user.CreateDB(u.conn); err != nil {
			log.Fatalf("Read database error")
		}
		return user, true
	}
	rows.Scan(user.Name, user.ChatID, user.Location)
	return user, false
}

// TODO: it strange arguments
func (u UsersDB) Create(id UserID, user User) {
	if err := user.CreateDB(u.conn); err != nil {
		log.Fatalf("Read database error")
	}
}
func (u UsersDB) Update(id UserID, user User) error {
	return user.UpdateDB(u.conn)
}
