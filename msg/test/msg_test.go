package test

import (
	"database/sql"
	"github.com/davecgh/go-spew/spew"
	_ "github.com/mattn/go-sqlite3"
	"os/user"
	"testing"
)

func TestReadMessage(t *testing.T) {
	t.Run("test check read message", func(t *testing.T) {
		user,err := user.Current()
		if err != nil {
			spew.Dump(err)
		}
		spew.Dump(user.HomeDir)
		db, openErr := sql.Open("sqlite3", "/Users/sunbowen/chat.db")
		if openErr != nil {
			spew.Dump(openErr)
		}
		rows, queryErr := db.Query("begin;\nselect text from message where ROWID IN (select max(ROWID) from message);")
		if queryErr != nil {
			spew.Dump(queryErr)
		}
		for rows.Next() {
			var text string
			rows.Scan(&text)
			spew.Dump(text)
		}
	})
}
