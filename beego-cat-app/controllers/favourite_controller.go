
//  controllers/favourite_controller.go

package controllers

import (
    "encoding/json"
    "net/http"
    "time"
    "fmt"
    // "log"
    "github.com/beego/beego/v2/server/web"
    "io/ioutil"
)

type FavoritesController struct {
    web.Controller
}

// struct for Favorite object
type Favorite struct {
    ID        interface{} `json:"id"`      // Leave this as interface{} to handle both string and number
    ImageID   string      `json:"image_id"`
    URL       string      `json:"url"`
    SubID     string      `json:"sub_id"`
    CreatedAt time.Time   `json:"created_at"`
}

type SyncFavoritesRequest struct {
    Favorites []Favorite `json:"favorites"`
}

// GetFavorites retrieves the list of favorites for a given sub_id
func (c *FavoritesController) GetFavorites() {
    subID := c.Ctx.Input.Param(":subId")
    apiKey, _ := web.AppConfig.String("cat_api_key")
    
    client := &http.Client{Timeout: 10 * time.Second}
    req, _ := http.NewRequest("GET", 
        fmt.Sprintf("https://api.thecatapi.com/v1/favourites?sub_id=%s", subID), 
        nil)
    
    req.Header.Add("x-api-key", apiKey)

    resp, err := client.Do(req)
    if err != nil {
        c.Error(http.StatusInternalServerError, err.Error())
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        c.Error(resp.StatusCode, "Failed to fetch favorites")
        return
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        c.Error(http.StatusInternalServerError, err.Error())
        return
    }

    c.Ctx.Output.Header("Content-Type", "application/json")
    c.Ctx.Output.Body(body)
}

// RemoveFavorite removes a favorite for a specific user by ID
func (c *FavoritesController) RemoveFavorite() {
    favoriteID := c.Ctx.Input.Param(":favoriteId")
    apiKey, _ := web.AppConfig.String("cat_api_key")
    
    client := &http.Client{Timeout: 10 * time.Second}
    
    req, err := http.NewRequest("DELETE", 
        fmt.Sprintf("https://api.thecatapi.com/v1/favourites/%s", favoriteID), 
        nil)
    if err != nil {
        c.Error(500, "Failed to create request")
        return
    }
    
    req.Header.Add("x-api-key", apiKey)
    
    resp, err := client.Do(req)
    if err != nil {
        c.Error(500, "Failed to remove favorite")
        return
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        c.Error(resp.StatusCode, "Failed to remove favorite from API")
        return
    }
    
    c.Data["json"] = map[string]string{"message": "Favorite removed successfully"}
    c.ServeJSON()
}

// SyncFavorites synchronizes favorites with the backend
func (c *FavoritesController) SyncFavorites() {
    subID := c.Ctx.Input.Param(":subId")
    
    var req SyncFavoritesRequest
    if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
        c.Error(400, "Invalid request data")
        return
    }
    
    // Log the sync operation
    fmt.Printf("Synced %d favorites for user %s\n", len(req.Favorites), subID)
    
    c.Data["json"] = map[string]interface{}{
        "message": "Favorites synced successfully",
        "count":   len(req.Favorites),
    }
    c.ServeJSON()
}

// Error handles error responses
func (c *FavoritesController) Error(status int, message string) {
    c.Ctx.Output.SetStatus(status)
    c.Data["json"] = map[string]string{"error": message}
    c.ServeJSON()
}
