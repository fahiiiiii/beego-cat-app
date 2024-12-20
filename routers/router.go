// routers/router.go
package routers

import (
	"beego-cat-app/controllers"
	"github.com/beego/beego/v2/server/web"
)

func init() {
	web.Router("/", &controllers.MainController{})
	web.Router("/favorites", &controllers.MainController{}, "get:ShowFavorites")
	web.Router("/breeds", &controllers.MainController{}, "get:ShowBreeds")
	web.Router("/voting", &controllers.MainController{}, "get:ShowVoting")
	web.Router("/api/cats", &controllers.CatController{}, "get:FetchCats")
	// web.Router("/voting", &web.Controller{}, "get:VotingPage")
	


	
}