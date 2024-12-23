// controllers/cat_controller.go
package controllers

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "time"
    "fmt"
    "github.com/beego/beego/v2/server/web"
    "beego-cat-app/models"
)

type CatController struct {
    web.Controller
}

// type CatImage struct {
//     URL string `json:"url"`
//     ID  string `json:"id"`
// }

type ImageChannel struct {
    images chan []models.CatImage
    errors chan error
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

// GetImages handles the GET request for fetching cat images
// func (c *CatController) GetImages() {
//     // Set CORS headers
//     c.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
//     c.Ctx.Output.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
//     c.Ctx.Output.Header("Access-Control-Allow-Headers", "Content-Type")

//     // Handle OPTIONS request
//     if c.Ctx.Request.Method == "OPTIONS" {
//         c.Ctx.Output.SetStatus(200)
//         return
//     }

//     imgChan := NewImageChannel()
    
//     // Fetch images asynchronously
//     go c.fetchCatImages(imgChan)

//     // Wait for response with timeout
//     select {
//     case images := <-imgChan.images:
//         c.Data["json"] = images
//         c.ServeJSON()
//     case err := <-imgChan.errors:
//         c.Ctx.Output.SetStatus(http.StatusInternalServerError)
//         c.Data["json"] = map[string]string{"error": err.Error()}
//         c.ServeJSON()
//     case <-time.After(12 * time.Second):
//         c.Ctx.Output.SetStatus(http.StatusRequestTimeout)
//         c.Data["json"] = map[string]string{"error": "Request timeout"}
//         c.ServeJSON()
//     }
// }
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
}