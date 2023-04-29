package models

import (
	"database/sql"
	"os"
	"testing"
)

func newTestDb(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "test_web:pass@/test_snippetbox?parseTime=true&multiStatements=true")

	if err != nil {
		t.Fatal(err)
	}

	script, err := os.ReadFile("./migrations/setup.sql")

	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(script))

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		script, err := os.ReadFile("./migrations/teardown.sql")

		if err != nil {
			t.Fatal(err)
		}

		_, err = db.Exec(string(script))

		if err != nil {
			t.Fatal(err)
		}

		db.Close()
	})

	return db
}
