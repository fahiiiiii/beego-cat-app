package controllers

import (
    "beego-cat-app/models"
    "fmt"
    "github.com/beego/beego/v2/server/web"
    "log"
    "net/http"
    "encoding/json"
    "sync"
)

// FavoriteController handles requests related to favorite cat images
type FavoriteController struct {
    web.Controller
}

// ShowFavorites fetches the favorite images and sends them as a JSON response
func (c *FavoriteController) ShowFavorites() {
    favoriteIDs := []string{"abc123", "def456", "ghi789"} // Example favorite IDs

    // Create a channel to manage concurrent fetching
    ch := make(chan models.CatImage)
    var wg sync.WaitGroup

    // Launch a goroutine for each favorite image ID
    for _, id := range favoriteIDs {
        wg.Add(1)
        go fetchCatImage(id, ch, &wg)
    }

    // Collect the responses
    var catImages []models.CatImage
    go func() {
        // Wait for all goroutines to finish
        wg.Wait()
        close(ch) // Close the channel after all goroutines are done
    }()

    // Collect cat images from the channel
    for catImage := range ch {
        catImages = append(catImages, catImage)
    }

    // Return the favorite cat images as a JSON response
    c.Data["json"] = catImages
    c.ServeJSON()
}

// Fetch a cat image from the CatAPI
func fetchCatImage(id string, ch chan<- models.CatImage, wg *sync.WaitGroup) {
    defer wg.Done() // Mark the goroutine as done when it finishes

    // URL to fetch the cat image details from the CatAPI
    apiURL := fmt.Sprintf("https://api.thecatapi.com/v1/images/%s", id)

    // Make the request to CatAPI
    req, err := http.NewRequest("GET", apiURL, nil)
    if err != nil {
        log.Println("Error creating request:", err)
        return
    }
	API_KEY, _ := web.AppConfig.String("cat_api_key")
    // BASE_URL, _ := web.AppConfig.String("cat_api_base_url")
    // Add the API key to the request header
    req.Header.Add("x-api-key", API_KEY)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Println("Error making request to CatAPI:", err)
        return
    }
    defer resp.Body.Close()

    // Check for a non-200 status code (indicating an error)
    if resp.StatusCode != http.StatusOK {
        log.Println("Error: Received non-200 response status code:", resp.StatusCode)
        return
    }

    // Parse the JSON response
    var catImage models.CatImage
    if err := json.NewDecoder(resp.Body).Decode(&catImage); err != nil {
        log.Println("Error decoding JSON response:", err)
        return
    }

    // Send the fetched image details to the channel
    ch <- catImage
}
