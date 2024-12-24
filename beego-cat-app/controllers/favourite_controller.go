// controllers/favourite_controller.go

package controllers

import (
    "encoding/json"
    "net/http"
    "time"
    "fmt"
    "github.com/beego/beego/v2/server/web"
)

type FavoritesController struct {
    web.Controller
}

type Favorite struct {
    ID        string    `json:"id"`
    ImageID   string    `json:"image_id"`
    URL       string    `json:"url"`
    SubID     string    `json:"sub_id"`
    CreatedAt time.Time `json:"created_at"`
}

type SyncFavoritesRequest struct {
    Favorites []Favorite `json:"favorites"`
}

// GetUserFavorites retrieves favorites for a specific user
func (c *FavoritesController) GetUserFavorites() {
    subID := c.Ctx.Input.Param(":subId")
    apiKey, _ := web.AppConfig.String("cat_api_key")
    
    // Create HTTP client with timeout
    client := &http.Client{Timeout: 10 * time.Second}
    
    // Make request to The Cat API
    req, err := http.NewRequest("GET", 
        fmt.Sprintf("https://api.thecatapi.com/v1/favourites?sub_id=%s", subID), 
        nil)
    if err != nil {
        c.Error(500, "Failed to create request")
        return
    }
    
    req.Header.Add("x-api-key", apiKey)
    
    resp, err := client.Do(req)
    if err != nil {
        c.Error(500, "Failed to fetch favorites")
        return
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        c.Error(resp.StatusCode, "Failed to get favorites from API")
        return
    }
    
    // Parse and return the favorites
    var favorites []Favorite
    if err := json.NewDecoder(resp.Body).Decode(&favorites); err != nil {
        c.Error(500, "Failed to parse favorites")
        return
    }
    
    c.Data["json"] = favorites
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
    
    // Here you could save the favorites to your database
    // For example:
    // err := models.SaveUserFavorites(subID, req.Favorites)
    
    // For now, we'll just log them
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