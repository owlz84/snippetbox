package models

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"os"
	"snippetbox.stuartlynn.net/internal/args"
	"testing"
)

func newTestDB(t *testing.T) *sql.DB {
	// Establish a sql.DB connection pool for our test database. Because our
	// setup and teardown scripts contains multiple SQL statements, we need
	// to use the "multiStatements=true" parameter in our DSN. This instructs
	// our MySQL database driver to support executing multiple SQL statements
	// in one db.Exec() call
	testArgs := args.ParseArgs("config_test.yml", "../../")

	var databaseConfig mysql.Config
	databaseConfig.User = testArgs.User
	databaseConfig.Passwd = testArgs.Pwd
	databaseConfig.Net = "tcp"
	databaseConfig.Addr = testArgs.Host
	databaseConfig.DBName = testArgs.DB
	databaseConfig.ParseTime = true
	databaseConfig.MultiStatements = true
	databaseConfig.AllowNativePasswords = true

	dsn := databaseConfig.FormatDSN()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}

	// Read the setup SQL script from file and execute the statements.
	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	// Use the t.Cleanup() to register a function *which will automatically be
	// called by Go when the current test (or subtest) which calls newTestDB()
	// has finished*. In this function we read and execute the teardown script,
	// and close the database connection pool.
	t.Cleanup(func() {
		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}

		db.Close()
	})

	// Return the database connection pool.
	return db
}
