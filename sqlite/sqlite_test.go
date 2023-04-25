package sqlite

import (
	"os"
	"testing"
)

func TestSqlite(t *testing.T) {
	db := New("sqlite.db")
	db.With("key string, value string, i integer").
		Key("key").
		CreateTable("key_pair")
	_, err := os.Stat("sqlite.db")
	if err != nil {
		t.Errorf("unable to create sqlite.db")
	}

	db.From("key_pair").Replace("key001", "val001", 5)
	res := db.From("key_pair").Where("key = ?", "key001").Select("i")
	if !res.Next() {
		t.Error("unable to fetch key")
	} else {
		var i int
		res.Scan(&i)
		if 5 != i {
			t.Errorf("expected 5, got %d", i)
		}
	}

	t.Cleanup(func() {
		err := os.Remove("sqlite.db")
		if err != nil {
			t.Errorf("unable to remove sqlite.db")
		}
	})
}
