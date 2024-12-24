// controllers/unified_controller.go
package controllers

import (
    "encoding/json"
    "net/http"
    "sync"
    "github.com/beego/beego/v2/server/web"
	"beego-cat-app/models"
    "time"
    "fmt"
)

type UnifiedController struct {
    web.Controller
}

type UnifiedContent struct {
    Breeds    []Breed    `json:"breeds"`
    Favorites []Favorite `json:"favorites"`
    CatImages []models.CatImage `json:"catImages"`
    Error     string     `json:"error,omitempty"`
}

// fetchAllContent handles concurrent fetching of all content
func (c *UnifiedController) fetchAllContent() UnifiedContent {
    var content UnifiedContent
    var wg sync.WaitGroup
    var mu sync.Mutex

    // Fetch breeds
    wg.Add(1)
    go func() {
        defer wg.Done()
        breeds, err := fetchBreeds()
        mu.Lock()
        if err != nil {
            content.Error += "Failed to fetch breeds; "
        } else {
            content.Breeds = breeds
        }
        mu.Unlock()
    }()

    // Fetch favorites
    wg.Add(1)
    go func() {
        defer wg.Done()
        subID := c.Ctx.Input.Param(":subId")
        if subID != "" {
            apiKey, _ := web.AppConfig.String("cat_api_key")
            client := &http.Client{Timeout: 10 * time.Second}
            req, err := http.NewRequest("GET", 
                fmt.Sprintf("https://api.thecatapi.com/v1/favourites?sub_id=%s", subID), 
                nil)
            if err == nil {
                req.Header.Add("x-api-key", apiKey)
                resp, err := client.Do(req)
                if err == nil {
                    defer resp.Body.Close()
                    var favorites []Favorite
                    if json.NewDecoder(resp.Body).Decode(&favorites) == nil {
                        mu.Lock()
                        content.Favorites = favorites
                        mu.Unlock()
                    }
                }
            }
        }
    }()

    // Fetch cat images
    wg.Add(1)
    go func() {
        defer wg.Done()
        apiKey, _ := web.AppConfig.String("cat_api_key")
        baseURL, _ := web.AppConfig.String("cat_api_base_url")
        
        client := &http.Client{Timeout: 10 * time.Second}
        req, err := http.NewRequest("GET", fmt.Sprintf("%s/images/search?limit=10", baseURL), nil)
        if err == nil {
            req.Header.Add("x-api-key", apiKey)
            resp, err := client.Do(req)
            if err == nil {
                defer resp.Body.Close()
                var images []models.CatImage
                if json.NewDecoder(resp.Body).Decode(&images) == nil {
                    mu.Lock()
                    content.CatImages = images
                    mu.Unlock()
                }
            }
        }
    }()

    wg.Wait()
    return content
}

// Common handler for all pages
func (c *UnifiedController) servePageWithContent(template string) {
    // Fetch all content concurrently
    content := c.fetchAllContent()
    
    // Add content to template data
    c.Data["Breeds"] = content.Breeds
    c.Data["Favorites"] = content.Favorites
    c.Data["CatImages"] = content.CatImages
    c.Data["Error"] = content.Error
    
    // Set template
    c.TplName = template
}



// For DEBUG--------------------------------------
func (c *UnifiedController) Debug() {
    content := c.fetchAllContent()
    c.Data["json"] = content
    c.ServeJSON()
}
// ----------------------------------------------



// Handler for voting page
func (c *UnifiedController) ShowVoting() {
    c.servePageWithContent("voting.html")
}

// Handler for breeds page
func (c *UnifiedController) ShowBreeds() {
    c.servePageWithContent("breeds.html")
}

// Handler for favorites page
func (c *UnifiedController) ShowFavorites() {
    c.servePageWithContent("favorites.html")
}