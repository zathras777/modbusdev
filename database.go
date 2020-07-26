package modbusdev

import (
	"database/sql"
	"fmt"
	"log"

	// Import postgresql access
	_ "github.com/lib/pq"
)

// DatabaseConnection Structure that holds details of database connection and query to be executed
type DatabaseConnection struct {
	Host     string
	Port     uint
	User     string
	Password string
	Name     string

	SSL bool

	Query  string
	Fields []DatabaseField

	db        *sql.DB
	statement *sql.Stmt
}

// DatabaseField Struct to allow for managing a relationship between a named field in the database
// and the corresponding element in the device Map() data
type DatabaseField struct {
	Name string
	Code int
}

// OpenDatabase Open a database connection and call Ping to verify the connection.
func (dbC *DatabaseConnection) OpenDatabase() error {
	// Establish connection to postgresql....
	db, err := sql.Open("postgres", dbC.getConnectionString())
	if err != nil {
		return err
	}
	dbC.db = db
	err = dbC.db.Ping()
	if err != nil {
		dbC.db.Close()
		return err
	}
	return nil
}

func (dbC DatabaseConnection) getConnectionString() string {
	connDetails := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		dbC.Host, dbC.Port, dbC.User, dbC.Password, dbC.Name)
	if dbC.SSL {
		connDetails += " sslmode=enable"
	} else {
		connDetails += " sslmode=disable"
	}
	return connDetails
}

// Execute Execute the stored query using supplied map of values
func (dbC DatabaseConnection) Execute(data map[int]Value) error {
	if dbC.statement == nil {
		names := ""
		placeholders := ""
		for i, fld := range dbC.Fields {
			names += ", " + fld.Name
			placeholders += fmt.Sprintf(",$%d", i+1)
		}
		qry := fmt.Sprintf(dbC.Query, names, placeholders)
		stmt, err := dbC.db.Prepare(qry)
		if err != nil {
			log.Printf("Error creating insertion statement")
			return err
		}
		dbC.statement = stmt
	}
	qryData := make([]interface{}, len(dbC.Fields))
	for i, fld := range dbC.Fields {
		val, ck := data[fld.Code]
		if !ck {
			return fmt.Errorf("Code %d [%s] not listed in supplied map data", fld.Code, fld.Name)
		}
		qryData[i] = val.Ieee32
	}
	_, err := dbC.statement.Exec(qryData...)
	if err != nil {
		log.Printf("Error writing sensor data to database: %s", err)
		return err
	}
	return nil
}

// Close Close the database connection, closing any open statements as well.
func (dbC DatabaseConnection) Close() {
	if dbC.db == nil {
		return
	}
	if dbC.statement != nil {
		dbC.statement.Close()
	}
	dbC.db.Close()
}
