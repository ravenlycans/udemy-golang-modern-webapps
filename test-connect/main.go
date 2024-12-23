package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
)

func main() {
	// connect to database.
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:Centaur994699@localhost:5432/test_connect")
	if err != nil {
		log.Fatalf("Unable to connect to database, err: %v\n", err)
	}
	defer conn.Close(context.Background())

	// test the connection.
	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to ping datbase, err: %v\n", err)
	}

	// Get rows from table
	err = getAllRows(conn)
	if err != nil {
		log.Fatalf("Unable to get all rows from database, err: %v\n", err)
	}

	// insert a row
	query := `insert into users (first_name, last_name) values ($1, $2)`
	_, err = conn.Exec(context.Background(), query, "Jack", "Brown")
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Println("Inserted a row with Jack Brown")

	// Get rows from table.
	err = getAllRows(conn)

	// Update a row
	stmt := `update users set first_name = $1 where id=$2`
	_, err = conn.Exec(context.Background(), stmt, "Jackie", 4)
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}

	log.Println("Updated 1 or more rows.")

	// Get rows from table
	err = getAllRows(conn)

	// Get one row from table by id
	query = `select id, first_name, last_name from users where id=$1`
	row := conn.QueryRow(context.Background(), query, 3)

	var firstName, lastName string
	var id int
	err = row.Scan(&id, &firstName, &lastName)
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}

	log.Printf("QueryRow returns: %v, %s, %s\n", id, firstName, lastName)

	log.Println("---------------")

	// Delete a row
	query = "delete from users where id=$1"
	_, err = conn.Exec(context.Background(), query, 4)
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}

	log.Printf("Deleted row %v\n", 4)

	// Get rows from table
	err = getAllRows(conn)
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}
}

func getAllRows(conn *pgx.Conn) error {
	rows, err := conn.Query(context.Background(), "select id,first_name,last_name from users")
	if err != nil {
		log.Fatalf("getAllRows: Cannot get rows from database, err: %v\n", err)
	}
	defer rows.Close()

	var firstName, lastName string
	var id int

	for rows.Next() {
		err := rows.Scan(&id, &firstName, &lastName)
		if err != nil {
			log.Printf("Cannot extract id, firstname and lastname from rows.\n")
		}
		fmt.Printf("Record is: %v, %s, %s\n", id, firstName, lastName)
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("Error scanning rows.\n")
	}

	fmt.Println("-----------")

	return nil
}
