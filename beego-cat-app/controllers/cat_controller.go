// controllers/cat_controller.go
package controllers

import (
    "encoding/json"
    "bytes"
    "io"
    "io/ioutil"
    "net/http"
    "time"
    // "log"
    "fmt"
    "github.com/beego/beego/v2/server/web"
    "beego-cat-app/models"
)

// SaveFavoriteRequest represents the structure of the request body

type SaveFavoriteRequest struct {
    ImageID  string `json:"image_id"`
    ImageURL string `json:"image_url"`
}
type CatController struct {
    web.Controller
}


type ImageChannel struct {
    images chan []models.CatImage
    errors chan error
}


type VoteData struct {
    ImageID string `json:"image_id"`
    Value   int    `json:"value"`
}



func NewImageChannel() *ImageChannel {
    return &ImageChannel{
        images: make(chan []models.CatImage),
        errors: make(chan error, 1),
    }
}

func (c *CatController) fetchCatImages(imgChan *ImageChannel) {
    apiURL := "https://api.thecatapi.com/v1/images/search?limit=10"
    
    client := &http.Client{
        Timeout: 10 * time.Second,
    }
    
    resp, err := client.Get(apiURL)
    if err != nil {
        imgChan.errors <- err
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        imgChan.errors <- err
        return
    }
    var images []models.CatImage
    if err := json.Unmarshal(body, &images); err != nil {
        imgChan.errors <- err
        return
    }

    imgChan.images <- images
}


// GetImages handles fetching cat images
func (c *CatController) GetImages() {
    API_KEY, _ := web.AppConfig.String("cat_api_key")
    BASE_URL, _ := web.AppConfig.String("cat_api_base_url")
    fmt.Println("API Key:", API_KEY)    // Debug log
    fmt.Println("Base URL:", BASE_URL) // Debug log

    if API_KEY == "" || BASE_URL == "" {
        c.Error(500, "API key or Base URL is not configured")
        return
    }

    client := &http.Client{Timeout: 10 * time.Second}
    req, _ := http.NewRequest("GET", fmt.Sprintf("%s/images/search?limit=10", BASE_URL), nil)
    req.Header.Add("x-api-key", API_KEY)

    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("HTTP request error:", err)
        c.Error(500, err.Error())
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        fmt.Println("API response error:", resp.Status)
        c.Error(resp.StatusCode, fmt.Sprintf("API returned status: %s", resp.Status))
        return
    }
    var images []models.CatImage
    // var images []CatImage
    if err := json.NewDecoder(resp.Body).Decode(&images); err != nil {
        fmt.Println("JSON decoding error:", err)
        c.Error(500, err.Error())
        return
    }

    c.Data["json"] = images
    c.ServeJSON()
}


// GetRandomImage handles the GET request for fetching a single random cat image
func (c *CatController) GetRandomImage() {
    imgChan := NewImageChannel()
    
    go c.fetchCatImages(imgChan)

    select {
    case images := <-imgChan.images:
        if len(images) > 0 {
            c.Data["json"] = images[0] // Return just the first image
        } else {
            c.Data["json"] = map[string]string{"error": "No images found"}
        }
    case err := <-imgChan.errors:
        c.Ctx.Output.SetStatus(http.StatusInternalServerError)
        c.Data["json"] = map[string]string{"error": err.Error()}
    case <-time.After(12 * time.Second):
        c.Ctx.Output.SetStatus(http.StatusRequestTimeout)
        c.Data["json"] = map[string]string{"error": "Request timeout"}
    }

    c.ServeJSON()
}






func (c *CatController) SaveFavorite() {
    subID := c.Ctx.Input.Param(":subId")
    
    var reqData SaveFavoriteRequest
    if err := json.Unmarshal(c.Ctx.Input.RequestBody, &reqData); err != nil {
        c.Error(http.StatusBadRequest, "Invalid request data")
        return
    }

    // Log the favorite in your system
    fmt.Printf("Saving favorite image for subId: %s, ImageID: %s, URL: %s\n", 
        subID, reqData.ImageID, reqData.ImageURL)

    // Here you could save to your own database if needed
    // For example: SaveFavoriteToDatabase(subID, reqData.ImageID, reqData.ImageURL)

    c.Data["json"] = map[string]string{
        "message": "Favorite logged successfully",
        "sub_id": subID,
        "image_id": reqData.ImageID,
    }
    c.ServeJSON()
}

// GetFavorites retrieves favorites for a specific user from The Cat API
func (c *CatController) GetFavorites() {
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

// DeleteFavorite removes a favorite from The Cat API
func (c *CatController) DeleteFavorite() {
    favoriteID := c.Ctx.Input.Param(":favoriteId")
    apiKey, _ := web.AppConfig.String("cat_api_key")
    
    client := &http.Client{Timeout: 10 * time.Second}
    req, _ := http.NewRequest("DELETE", 
        fmt.Sprintf("https://api.thecatapi.com/v1/favourites/%s", favoriteID), 
        nil)
    
    req.Header.Add("x-api-key", apiKey)

    resp, err := client.Do(req)
    if err != nil {
        c.Error(http.StatusInternalServerError, err.Error())
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        c.Error(resp.StatusCode, "Failed to delete favorite")
        return
    }

    c.Data["json"] = map[string]string{"message": "Favorite deleted successfully"}
    c.ServeJSON()
}





func (c *CatController) SaveVote() {
    subID := c.Ctx.Input.Param(":subId")
    apiKey, _ := web.AppConfig.String("cat_api_key")
    
    var voteData VoteData
    if err := json.Unmarshal(c.Ctx.Input.RequestBody, &voteData); err != nil {
        c.Ctx.Output.SetStatus(http.StatusBadRequest)
        c.Data["json"] = map[string]string{"error": "Invalid request data"}
        c.ServeJSON()
        return
    }

    // Validate required fields
    if voteData.ImageID == "" || (voteData.Value != 0 && voteData.Value != 1) {
        c.Ctx.Output.SetStatus(http.StatusBadRequest)
        c.Data["json"] = map[string]string{"error": "Invalid vote data"}
        c.ServeJSON()
        return
    }

    // Prepare request to The Cat API
    payload := map[string]interface{}{
        "image_id": voteData.ImageID,
        "sub_id":   subID,
        "value":    voteData.Value,
    }

    jsonData, _ := json.Marshal(payload)
    
    // Send request to The Cat API
    client := &http.Client{Timeout: 10 * time.Second}
    req, _ := http.NewRequest("POST", "https://api.thecatapi.com/v1/votes", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("x-api-key", apiKey)

    resp, err := client.Do(req)
    if err != nil {
        c.Ctx.Output.SetStatus(http.StatusInternalServerError)
        c.Data["json"] = map[string]string{"error": "Failed to communicate with The Cat API"}
        c.ServeJSON()
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        body, _ := io.ReadAll(resp.Body)
        c.Ctx.Output.SetStatus(resp.StatusCode)
        c.Data["json"] = map[string]string{
            "error": fmt.Sprintf("The Cat API error: %s", string(body)),
        }
        c.ServeJSON()
        return
    }

    // Return success response
    c.Data["json"] = map[string]string{
        "message":  "Vote saved successfully",
        "sub_id":   subID,
        "image_id": voteData.ImageID,
    }
    c.ServeJSON()
}


// Prepare runs before each action
func (c *CatController) Prepare() {
    // Common setup code for all actions
    c.Ctx.Output.Header("Content-Type", "application/json")
}

// Error handles common error responses
func (c *CatController) Error(status int, message string) {
    c.Ctx.Output.SetStatus(status)
    c.Data["json"] = map[string]string{"error": message}
    c.ServeJSON()
}// SaveFavorite handles saving a favorite cat image to the server


