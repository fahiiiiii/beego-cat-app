// controllers/concurrent_controller.go
package controllers

import (
    "encoding/json"
    "net/http"
    "sync"
    "time"
    "github.com/beego/beego/v2/server/web"
	"fmt"
	"beego-cat-app/models"

)

// ConcurrentController combines all cat-related operations
type ConcurrentController struct {
    web.Controller
}

// Combined response structure
type CombinedResponse struct {
    Images    []models.CatImage `json:"images"`
    Breeds    []Breed          `json:"breeds"`
    Favorites []Favorite       `json:"favorites"`
    Votes     []VoteData       `json:"votes"`
    Errors    []string         `json:"errors,omitempty"`
}

// Channel structures for concurrent operations
type DataChannels struct {
    images    chan []models.CatImage
    breeds    chan []Breed
    favorites chan []Favorite
    votes     chan []VoteData
    errors    chan error
}

func NewDataChannels() *DataChannels {
    return &DataChannels{
        images:    make(chan []models.CatImage, 1),
        breeds:    make(chan []Breed, 1),
        favorites: make(chan []Favorite, 1),
        votes:     make(chan []VoteData, 1),
        errors:    make(chan error, 4), // Buffer for multiple potential errors
    }
}

// GetAllData handles concurrent fetching of all cat-related data
func (c *ConcurrentController) GetAllData() {
    subID := c.Ctx.Input.Param(":subId")
    apiKey, _ := web.AppConfig.String("cat_api_key")
    
    channels := NewDataChannels()
    var wg sync.WaitGroup
    
    // Start concurrent fetches
    wg.Add(4)
    
    // Fetch Images
    go func() {
        defer wg.Done()
        images, err := c.fetchImages(apiKey)
        if err != nil {
            channels.errors <- err
            return
        }
        channels.images <- images
    }()
    
    // Fetch Breeds
    go func() {
        defer wg.Done()
        breeds, err := fetchBreeds()
        if err != nil {
            channels.errors <- err
            return
        }
        channels.breeds <- breeds
    }()
    
    // Fetch Favorites
    go func() {
        defer wg.Done()
        favorites, err := c.fetchFavorites(apiKey, subID)
        if err != nil {
            channels.errors <- err
            return
        }
        channels.favorites <- favorites
    }()
    
    // Fetch Votes
    go func() {
        defer wg.Done()
        votes, err := c.fetchVotes(apiKey, subID)
        if err != nil {
            channels.errors <- err
            return
        }
        channels.votes <- votes
    }()
    
    // Wait for all goroutines in separate goroutine
    go func() {
        wg.Wait()
        close(channels.images)
        close(channels.breeds)
        close(channels.favorites)
        close(channels.votes)
        close(channels.errors)
    }()
    
    // Collect results with timeout
    response := CombinedResponse{}
    errors := []string{}
    
    timeout := time.After(15 * time.Second)
    
    // Collect data from channels
    for i := 0; i < 4; i++ {
        select {
        case imgs := <-channels.images:
            response.Images = imgs
        case breeds := <-channels.breeds:
            response.Breeds = breeds
        case favs := <-channels.favorites:
            response.Favorites = favs
        case votes := <-channels.votes:
            response.Votes = votes
        case err := <-channels.errors:
            if err != nil {
                errors = append(errors, err.Error())
            }
        case <-timeout:
            errors = append(errors, "Operation timed out")
            goto TimeoutOccurred
        }
    }
    
TimeoutOccurred:
    if len(errors) > 0 {
        response.Errors = errors
    }
    
    c.Data["json"] = response
    c.ServeJSON()
}

// Helper methods for fetching different types of data
func (c *ConcurrentController) fetchImages(apiKey string) ([]models.CatImage, error) {
    client := &http.Client{Timeout: 10 * time.Second}
    req, _ := http.NewRequest("GET", "https://api.thecatapi.com/v1/images/search?limit=10", nil)
    req.Header.Add("x-api-key", apiKey)
    
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var images []models.CatImage
    if err := json.NewDecoder(resp.Body).Decode(&images); err != nil {
        return nil, err
    }
    return images, nil
}

func (c *ConcurrentController) fetchFavorites(apiKey, subID string) ([]Favorite, error) {
    client := &http.Client{Timeout: 10 * time.Second}
    req, _ := http.NewRequest("GET", 
        fmt.Sprintf("https://api.thecatapi.com/v1/favourites?sub_id=%s", subID), 
        nil)
    req.Header.Add("x-api-key", apiKey)
    
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var favorites []Favorite
    if err := json.NewDecoder(resp.Body).Decode(&favorites); err != nil {
        return nil, err
    }
    return favorites, nil
}

func (c *ConcurrentController) fetchVotes(apiKey, subID string) ([]VoteData, error) {
    client := &http.Client{Timeout: 10 * time.Second}
    req, _ := http.NewRequest("GET", 
        fmt.Sprintf("https://api.thecatapi.com/v1/votes?sub_id=%s", subID), 
        nil)
    req.Header.Add("x-api-key", apiKey)
    
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var votes []VoteData
    if err := json.NewDecoder(resp.Body).Decode(&votes); err != nil {
        return nil, err
    }
    return votes, nil
}