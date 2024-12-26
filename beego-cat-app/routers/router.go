// routers/router.go
package routers

import (
    "beego-cat-app/controllers" // Ensure this path is correct
    "github.com/beego/beego/v2/server/web"
)

func init() {
    // UnifiedController - if you have such a controller, ensure it's imported correctly
    // unified := &controllers.UnifiedController{}
    // web.Router("/", unified, "get:ShowVoting")
    // MainController routes
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
	
	// Routing setup in Go (e.g., in `routers/router.go`)
	// web.Router("/favorites/:subId", &controllers.FavoritesController{}, "get:GetUserFavorites")
	web.Router("/favorites/:subId", &controllers.FavoritesController{}, "get:GetFavorites")
	
	// Add this to router.go
	// web.Router("/api/vote", &controllers.CatController{}, "post:AddVote")
	// routers/router.go
	// Add this to router.go
	

	// web.Router("/api/vote/:subId", &controllers.CatController{}, "post:SaveVote")
	web.Router("/api/vote/:subId", &controllers.CatController{}, "post:SaveVote")
	 // New concurrent data loading endpoints
    web.Router("/api/concurrent/voting-data", &controllers.UnifiedController{}, "get:GetVotingData")
    web.Router("/api/concurrent/favorites-data", &controllers.UnifiedController{}, "get:GetFavoritesData")
    web.Router("/api/concurrent/breeds-data", &controllers.UnifiedController{}, "get:GetBreedsData")

	// web.Router("/api/vote/:subId", &controllers.CatController{}, "post:SaveVote")
	// web.Router("/api/vote/:user_id", &controllers.VoteController{}, "post:Post")

	
}
// // routers/router.go
// package routers

// import (
//     "beego-cat-app/controllers"
//     "github.com/beego/beego/v2/server/web"
// )

// func init() {
//     // Existing routes
//     web.Router("/", &controllers.MainController{})
//     web.Router("/favorites", &controllers.MainController{}, "get:ShowFavorites")
//     web.Router("/breeds", &controllers.MainController{}, "get:ShowBreeds")
//     web.Router("/voting", &controllers.MainController{}, "get:ShowVoting")
    
//     // Existing API routes
//     web.Router("/api/getImages", &controllers.CatController{}, "get:GetImages")
//     web.Router("/api/cat/random-image", &controllers.CatController{}, "get:GetRandomImage")
//     web.Router("/api/cat/save-favorite/:subId", &controllers.CatController{}, "post:SaveFavorite")
//     web.Router("/api/cat/favorites/:subId", &controllers.CatController{}, "get:GetFavorites")
//     web.Router("/api/cat/delete-favorite/:favoriteId", &controllers.CatController{}, "delete:DeleteFavorite")
//     web.Router("/api/breeds", &controllers.BreedController{}, "get:GetBreeds")
//     web.Router("/api/breed-images", &controllers.BreedController{}, "get:GetBreedImages")
//     web.Router("/favorites/:subId", &controllers.FavoritesController{}, "get:GetFavorites")
//     web.Router("/api/vote/:subId", &controllers.CatController{}, "post:SaveVote")
    
//     // New concurrent data loading endpoints
//     web.Router("/api/concurrent/voting-data", &controllers.UnifiedController{}, "get:GetVotingData")
//     web.Router("/api/concurrent/favorites-data", &controllers.UnifiedController{}, "get:GetFavoritesData")
//     web.Router("/api/concurrent/breeds-data", &controllers.UnifiedController{}, "get:GetBreedsData")
// }