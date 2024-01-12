package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	//_ "github.com/lib/pq"

	_ "github.com/stretchr/testify/require"
	_ "gorm.io/driver/postgres"
	_ "gorm.io/gorm"
)

var testQueries *Queries
var testDB *sql.DB

const (
	dbDriver = "pgx"
	dbSource = "postgresql://root:mysecret@localhost:5433/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal(err)
	}
	testQueries = New(testDB)
	//testQueries.db.PrepareContext()
	os.Exit(m.Run())

}
