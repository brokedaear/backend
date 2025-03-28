package dal

import (
	"database/sql"
	"os"
	"testing"
)

func newTestDB(t *testing.T) *sql.DB {
	// Open DB connection

	// Read the SQL script

	script, err := os.ReadFile("location/of/setup.sql")
	if err != nil {
		// close db
		t.Fatal(err)
	}

	// Execute the script query

	t.Cleanup(func() {
		// defer db.close

		script, err := os.ReadFile("localtion/of/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		// execute script query
	})

	// return the db
}
