package controllers

import (
	"time"

	"github.com/revel/revel"
)

var start = time.Now()

// App is the main controller
type App struct {
	*revel.Controller
}

// Index action of App controller
func (c App) Index() revel.Result {

	moreScripts := moreScripts()
	return c.Render(moreScripts)
}

func moreScripts() []string {
	result := make([]string, 0)
	result = append(result, "js/app/app.js")
	result = append(result, "js/app/ctrl/HomeCtrl.js")

	return result
}
