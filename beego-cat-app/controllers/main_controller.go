//controllers/main_controller.go

package controllers

import (
	"github.com/beego/beego/v2/server/web"
)

type MainController struct {
	web.Controller
}

// Default handler for the root URL ("/")
func (c *MainController) Get() {
	c.TplName = "voting.html" // Render the voting.html template by default
}

// ShowFavorites renders the favorites page
func (c *MainController) ShowFavorites() {
	c.TplName = "favorites.html"
}

// ShowBreeds renders the breeds page
func (c *MainController) ShowBreeds() {
	c.TplName = "breeds.html"
}

// ShowVoting renders the voting page
func (c *MainController) ShowVoting() {
	c.TplName = "voting.html"
}






// For DEBUG
func (c *UnifiedController) ShowDebug() {
    c.TplName = "debug.html"
}