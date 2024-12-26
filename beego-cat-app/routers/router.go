// routers/router.go
package routers

import (
    "beego-cat-app/controllers" // Ensure this path is correct
    "github.com/beego/beego/v2/server/web"
)

func init() {
    
    web.Router("/", &controllers.MainController{})
    web.Router("/favorites", &controllers.MainController{}, "get:ShowFavorites")
    web.Router("/breeds", &controllers.MainController{}, "get:ShowBreeds")
    web.Router("/voting", &controllers.MainController{}, "get:ShowVoting")

    // CatController routes
    web.Router("/api/getImages", &controllers.CatController{}, "get:GetImages")
    web.Router("/api/cat/random-image", &controllers.CatController{}, "get:GetRandomImage")
	web.Router("/api/cat/save-favorite/:subId", &controllers.CatController{}, "post:SaveFavorite")
    web.Router("/api/cat/favorites/:subId", &controllers.CatController{}, "get:GetFavorites")
    web.Router("/api/cat/delete-favorite/:favoriteId", &controllers.CatController{}, "delete:DeleteFavorite")

	
	web.Router("/api/breeds", &controllers.BreedController{}, "get:GetBreeds")
	web.Router("/api/breed-images", &controllers.BreedController{}, "get:GetBreedImages")
	
	web.Router("/favorites/:subId", &controllers.FavoritesController{}, "get:GetFavorites")
	
	web.Router("/api/vote/:subId", &controllers.CatController{}, "post:SaveVote")
	 // New concurrent data loading endpoints
    web.Router("/api/concurrent/voting-data", &controllers.UnifiedController{}, "get:GetVotingData")
    web.Router("/api/concurrent/favorites-data", &controllers.UnifiedController{}, "get:GetFavoritesData")
    web.Router("/api/concurrent/breeds-data", &controllers.UnifiedController{}, "get:GetBreedsData")

	
	
}