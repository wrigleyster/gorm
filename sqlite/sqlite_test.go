package sqlite

import (
	"os"
	"testing"
)

func TestSqlite(t *testing.T) {
	db := New("sqlite.db")
	db.With("key string, value string").
		Key("key").
		CreateTable("key_pair")
	_, err := os.Stat("sqlite.db")
	if err != nil {
		t.Errorf("unable to create sqlite.db")
	}
	t.Cleanup(func() {
		err := os.Remove("sqlite.db")
		if err != nil {
			t.Errorf("unable to remove sqlite.db")
		}
	})
}
