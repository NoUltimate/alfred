package main

import (
	"alfred/util"
	"database/sql"
	aw "github.com/deanishe/awgo"
	_ "github.com/mattn/go-sqlite3"
	"os/user"
	"strings"
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
	user,err := user.Current()
	if err != nil {
		wf.NewItem("获取home路径失败")
		wf.SendFeedback()
	}
	db, openErr := sql.Open("sqlite3", user.HomeDir + "/Library/Messages/chat.db")
	if openErr != nil {
		wf.NewItem("读取数据库失败")
		wf.SendFeedback()
	}
	rows, queryErr := db.Query("begin;\nselect text from message where ROWID IN (select ROWID from message order by ROWID desc limit 0,10) order by ROWID desc;")
	if queryErr != nil {
		wf.NewItem("读取数据失败")
		wf.SendFeedback()
	}
	for rows.Next() {
		var text string
		rows.Scan(&text)
		if !strings.Contains(text, "验证码") {
			continue
		}
		code := util.CheckMessageCode(text, 6)
		if len(code) == 0 {
			code = util.CheckMessageCode(text, 4)
		}
		wf.NewItem(code).Subtitle(text).Arg(code).Valid(true)
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
