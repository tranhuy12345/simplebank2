package db

import (
	"database/sql"
	"db/db/util"
	"log"
	"os"
	"testing"

	//_ "github.com/lib/pq"

	_ "github.com/lib/pq"
	_ "github.com/stretchr/testify/require"
	_ "gorm.io/driver/postgres"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var config util.Config
	var err error
	config, err = util.LoadConfig("../..")
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}
	testQueries = New(testDB)
	//testQueries.db.PrepareContext()
	os.Exit(m.Run())

}
