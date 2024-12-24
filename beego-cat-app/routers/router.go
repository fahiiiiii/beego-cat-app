// routers/router.go
package routers


import (
	"beego-cat-app/controllers"
	"github.com/beego/beego/v2/server/web"
    

	
)

func init() {

	unified := &controllers.UnifiedController{}
    web.Router("/", unified, "get:ShowVoting")
    // web.Router("/voting", unified, "get:ShowVoting")
    // web.Router("/breeds", unified, "get:ShowBreeds")
    // web.Router("/favorites", unified, "get:ShowFavorites")

	web.Router("/", &controllers.MainController{})
	web.Router("/favorites", &controllers.MainController{}, "get:ShowFavorites")
	web.Router("/breeds", &controllers.MainController{}, "get:ShowBreeds")
	web.Router("/voting", &controllers.MainController{}, "get:ShowVoting")
	
	// API routes for cat breeds and breed images
	web.Router("/api/breeds", &controllers.BreedController{}, "get:GetBreeds") // Fetch breeds
	web.Router("/api/breed-images", &controllers.BreedController{}, "get:GetBreedImages") // Fetch breed images


	web.Router("/api/breeds", &controllers.BreedController{}, "get:GetBreeds")
    web.Router("/api/breed-images", &controllers.BreedController{}, "get:GetBreedImages")
	
	
	web.Router("/api/favorites/:subId", &controllers.FavoritesController{}, "get:GetUserFavorites")
	web.Router("/api/syncFavorites/:subId", &controllers.FavoritesController{}, "post:SyncFavorites")


	// For DEBUG-----------------------------
	web.Router("/debug", unified, "get:Debug")
	web.Router("/debug-view", unified, "get:ShowDebug")
	
	
}