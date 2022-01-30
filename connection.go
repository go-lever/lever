package lever

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

const (
	defaultTimeout = 10*time.Second
)

// ConnectionParams represents the database connection parameters
// The following parameters are needed : Host, Port, User, Password, DBName
type ConnectionParams struct {
	Host string
	Port string
	User string
	Password string
	DBName string
	TimeOut time.Duration
}

// NewPostgresConnection creates a new sql.DB compatible Postgres connection
// After having opened the connection, NewPostgresConnection tries to ping the database
// and wait for a timeout.
//
// NewPostgresConnection uses *ConnectionParams to open the connection.
func NewPostgresConnection(param *ConnectionParams) *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	param.Host, param.Port, param.User, param.Password, param.DBName)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatalf("unable to open database connection : %s", err.Error())
	}

	to := defaultTimeout
	if param.TimeOut.Milliseconds() == 0 {
		to = param.TimeOut
	}
	
	err = db.Ping()
	if err != nil {
		log.Println("waiting for DB to be ready")
		//wait for DB connection
		time.Sleep(to)
		if errS := db.Ping(); errS != nil {
			log.Fatalf("unable to ping database : %s", err.Error())
		}
	}

	return db
}