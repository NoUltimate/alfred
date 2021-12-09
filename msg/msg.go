package main

import (
	"database/sql"
	aw "github.com/deanishe/awgo"
	_ "github.com/mattn/go-sqlite3"
	"os/user"
)

// Workflow is the main API
var wf *aw.Workflow

func init() {
	// Create a new Workflow using default settings.
	// Critical settings are provided by Alfred via environment variables,
	// so this *will* die in flames if not run in an Alfred-like environment.
	wf = aw.New()
}

// Your workflow starts here
func run() {
	items := make([]string, 0)
	user,err := user.Current()
	if err != nil {
		items = append(items, "获取home路径失败")
	}
	db, openErr := sql.Open("sqlite3", user.HomeDir + "/Library/Messages/chat.db")
	if openErr != nil {
		items = append(items, "读取数据库失败")
	}

	args := wf.Args()
	limit := "50"
	query := "\"%" + args[0] + "%\""

	if len(args) > 1 {
		limit = args[1]
	}
	rows, queryErr := db.Query("begin;\nselect text from message where ROWID IN " +
		"(select ROWID from message order by ROWID desc limit 0," + limit + ")" +
		"and text like " + query + " order by ROWID desc;")
	if queryErr != nil {
		items = append(items, "读取数据失败")
	}
	if rows != nil {
		for rows.Next() {
			var text string
			rows.Scan(&text)
			items = append(items, text)
		}
	}

	for _, item := range items {
		wf.NewItem(item).Arg(item).Valid(true)
	}
	wf.SendFeedback()
	// Add a "Script Filter" result
	// Send results to Alfred
}

func main() {
	// Wrap your entry point with Run() to catch and log panics and
	// show an error in Alfred instead of silently dying
	wf.Run(run)
}
