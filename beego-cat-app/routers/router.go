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
	// web.Router("/api/cats", &controllers.CatController{}, "get:FetchCats")
	// web.Router("/breeds", &controllers.BreedController{}, "get:FetchBreeds")

	// beego.Router("/breeds", &controllers.BreedController{}, "get:GetBreeds")


	// API routes for cat breeds and breed images
	web.Router("/api/breeds", &controllers.BreedController{}, "get:GetBreeds") // Fetch breeds
	web.Router("/api/breed-images", &controllers.BreedController{}, "get:GetBreedImages") // Fetch breed images


	web.Router("/api/breeds", &controllers.BreedController{}, "get:GetBreeds")
    web.Router("/api/breed-images", &controllers.BreedController{}, "get:GetBreedImages")
	
	
	web.Router("/favorites", &controllers.FavoriteController{}, "get:ShowFavorites") // Fetch and display favorites
}

