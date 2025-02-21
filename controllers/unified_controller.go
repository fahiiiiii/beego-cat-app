// // controllers/unified_controller.go
// package controllers

// import (
//     "encoding/json"
//     "net/http"
//     "sync"
//     "github.com/beego/beego/v2/server/web"
// 	"beego-cat-app/models"
//     "time"
//     "fmt"
// )

// type UnifiedController struct {
//     web.Controller
// }

// type UnifiedContent struct {
//     Breeds    []Breed    `json:"breeds"`
//     Favorites []Favorite `json:"favorites"`
//     CatImages []models.CatImage `json:"catImages"`
//     Error     string     `json:"error,omitempty"`
// }

// // fetchAllContent handles concurrent fetching of all content
// func (c *UnifiedController) fetchAllContent() UnifiedContent {
//     var content UnifiedContent
//     var wg sync.WaitGroup
//     var mu sync.Mutex

//     // Fetch breeds
//     wg.Add(1)
//     go func() {
//         defer wg.Done()
//         breeds, err := fetchBreeds()
//         mu.Lock()
//         if err != nil {
//             content.Error += "Failed to fetch breeds; "
//         } else {
//             content.Breeds = breeds
//         }
//         mu.Unlock()
//     }()

//     // Fetch favorites
//     wg.Add(1)
//     go func() {
//         defer wg.Done()
//         subID := c.Ctx.Input.Param(":subId")
//         if subID != "" {
//             apiKey, _ := web.AppConfig.String("cat_api_key")
//             client := &http.Client{Timeout: 10 * time.Second}
//             req, err := http.NewRequest("GET", 
//                 fmt.Sprintf("https://api.thecatapi.com/v1/favourites?sub_id=%s", subID), 
//                 nil)
//             if err == nil {
//                 req.Header.Add("x-api-key", apiKey)
//                 resp, err := client.Do(req)
//                 if err == nil {
//                     defer resp.Body.Close()
//                     var favorites []Favorite
//                     if json.NewDecoder(resp.Body).Decode(&favorites) == nil {
//                         mu.Lock()
//                         content.Favorites = favorites
//                         mu.Unlock()
//                     }
//                 }
//             }
//         }
//     }()

//     // Fetch cat images
//     wg.Add(1)
//     go func() {
//         defer wg.Done()
//         apiKey, _ := web.AppConfig.String("cat_api_key")
//         baseURL, _ := web.AppConfig.String("cat_api_base_url")
        
//         client := &http.Client{Timeout: 10 * time.Second}
//         req, err := http.NewRequest("GET", fmt.Sprintf("%s/images/search?limit=10", baseURL), nil)
//         if err == nil {
//             req.Header.Add("x-api-key", apiKey)
//             resp, err := client.Do(req)
//             if err == nil {
//                 defer resp.Body.Close()
//                 var images []models.CatImage
//                 if json.NewDecoder(resp.Body).Decode(&images) == nil {
//                     mu.Lock()
//                     content.CatImages = images
//                     mu.Unlock()
//                 }
//             }
//         }
//     }()

//     wg.Wait()
//     return content
// }

// // Common handler for all pages
// func (c *UnifiedController) servePageWithContent(template string) {
//     // Fetch all content concurrently
//     content := c.fetchAllContent()
    
//     // Add content to template data
//     c.Data["Breeds"] = content.Breeds
//     c.Data["Favorites"] = content.Favorites
//     c.Data["CatImages"] = content.CatImages
//     c.Data["Error"] = content.Error
    
//     // Set template
//     c.TplName = template
// }



// // For DEBUG--------------------------------------
// func (c *UnifiedController) Debug() {
//     content := c.fetchAllContent()
//     c.Data["json"] = content
//     c.ServeJSON()
// }
// // ----------------------------------------------



// // Handler for voting page
// func (c *UnifiedController) ShowVoting() {
//     c.servePageWithContent("voting.html")
// }

// // Handler for breeds page
// func (c *UnifiedController) ShowBreeds() {
//     c.servePageWithContent("breeds.html")
// }

// // Handler for favorites page
// func (c *UnifiedController) ShowFavorites() {
//     c.servePageWithContent("favorites.html")
// }
// controllers/unified_controller.go
package controllers

import (
    "encoding/json"
    "net/http"
    "sync"
    "time"
    "github.com/beego/beego/v2/server/web"
    "beego-cat-app/models"
    
    "fmt"
)

// UnifiedController handles concurrent data loading for each view
type UnifiedController struct {
    web.Controller
}

// Response structures for different views
type VotingViewData struct {
    RandomImage *models.CatImage `json:"randomImage"`
    Breeds     []Breed          `json:"breeds,omitempty"`
    UserVotes  []VoteData       `json:"userVotes,omitempty"`
    Errors     []string         `json:"errors,omitempty"`
}

type FavoritesViewData struct {
    Favorites []Favorite `json:"favorites"`
    Images   []models.CatImage `json:"images,omitempty"`
    Errors   []string    `json:"errors,omitempty"`
}

type BreedsViewData struct {
    Breeds []Breed  `json:"breeds"`
    Images []Image  `json:"images,omitempty"`
    Errors []string `json:"errors,omitempty"`
}

// GetVotingData handles concurrent data loading for voting view
func (c *UnifiedController) GetVotingData() {
    var wg sync.WaitGroup
    var mutex sync.Mutex
    data := VotingViewData{}
    errors := []string{}
    
    apiKey, _ := web.AppConfig.String("cat_api_key")
    subID := c.GetString("subId", "default-user") // Get from query or default
    
    // Channel for random image
    imageChan := make(chan *models.CatImage, 1)
    breedsChan := make(chan []Breed, 1)
    votesChan := make(chan []VoteData, 1)
    
    // Fetch random image
    wg.Add(1)
    go func() {
        defer wg.Done()
        client := &http.Client{Timeout: 10 * time.Second}
        req, _ := http.NewRequest("GET", "https://api.thecatapi.com/v1/images/search?limit=1", nil)
        req.Header.Add("x-api-key", apiKey)
        
        resp, err := client.Do(req)
        if err != nil {
            mutex.Lock()
            errors = append(errors, "Failed to fetch random image: "+err.Error())
            mutex.Unlock()
            return
        }
        defer resp.Body.Close()
        
        var images []models.CatImage
        if err := json.NewDecoder(resp.Body).Decode(&images); err != nil {
            mutex.Lock()
            errors = append(errors, "Failed to decode random image: "+err.Error())
            mutex.Unlock()
            return
        }
        
        if len(images) > 0 {
            imageChan <- &images[0]
        }
    }()
    
    // Fetch breeds concurrently
    wg.Add(1)
    go func() {
        defer wg.Done()
        breeds, err := fetchBreeds()
        if err != nil {
            mutex.Lock()
            errors = append(errors, "Failed to fetch breeds: "+err.Error())
            mutex.Unlock()
            return
        }
        breedsChan <- breeds
    }()
    
    // Fetch user votes concurrently
    wg.Add(1)
    go func() {
        defer wg.Done()
        client := &http.Client{Timeout: 10 * time.Second}
        req, _ := http.NewRequest("GET", 
            fmt.Sprintf("https://api.thecatapi.com/v1/votes?sub_id=%s", subID), 
            nil)
        req.Header.Add("x-api-key", apiKey)
        
        resp, err := client.Do(req)
        if err != nil {
            mutex.Lock()
            errors = append(errors, "Failed to fetch votes: "+err.Error())
            mutex.Unlock()
            return
        }
        defer resp.Body.Close()
        
        var votes []VoteData
        if err := json.NewDecoder(resp.Body).Decode(&votes); err != nil {
            mutex.Lock()
            errors = append(errors, "Failed to decode votes: "+err.Error())
            mutex.Unlock()
            return
        }
        votesChan <- votes
    }()
    
    // Wait for all goroutines with timeout
    done := make(chan bool)
    go func() {
        wg.Wait()
        done <- true
    }()
    
    // Handle timeout
    select {
    case <-done:
        // Collect results
        select {
        case data.RandomImage = <-imageChan:
        default:
        }
        select {
        case data.Breeds = <-breedsChan:
        default:
        }
        select {
        case data.UserVotes = <-votesChan:
        default:
        }
    case <-time.After(12 * time.Second):
        errors = append(errors, "Operation timed out")
    }
    
    if len(errors) > 0 {
        data.Errors = errors
    }
    
    c.Data["json"] = data
    c.ServeJSON()
}

// GetFavoritesData handles concurrent data loading for favorites view
func (c *UnifiedController) GetFavoritesData() {
    subID := c.GetString("subId", "default-user")
    apiKey, _ := web.AppConfig.String("cat_api_key")
    
    var wg sync.WaitGroup
    var mutex sync.Mutex
    data := FavoritesViewData{}
    errors := []string{}
    
    favoritesChan := make(chan []Favorite, 1)
    imagesChan := make(chan []models.CatImage, 1)
    
    // Fetch favorites
    wg.Add(1)
    go func() {
        defer wg.Done()
        favorites, err := c.fetchFavorites(apiKey, subID)  // Updated to use method receiver
        if err != nil {
            mutex.Lock()
            errors = append(errors, "Failed to fetch favorites: "+err.Error())
            mutex.Unlock()
            return
        }
        favoritesChan <- favorites
    }()
    
    // Fetch additional images
    wg.Add(1)
    go func() {
        defer wg.Done()
        images, err := c.fetchImages(apiKey)  // Updated to use method receiver
        if err != nil {
            mutex.Lock()
            errors = append(errors, "Failed to fetch images: "+err.Error())
            mutex.Unlock()
            return
        }
        imagesChan <- images
    }()
    
    // Handle completion or timeout
    done := make(chan bool)
    go func() {
        wg.Wait()
        done <- true
    }()
    
    select {
    case <-done:
        select {
        case data.Favorites = <-favoritesChan:
        default:
        }
        select {
        case data.Images = <-imagesChan:
        default:
        }
    case <-time.After(12 * time.Second):
        errors = append(errors, "Operation timed out")
    }
    
    if len(errors) > 0 {
        data.Errors = errors
    }
    
    c.Data["json"] = data
    c.ServeJSON()
}

// GetBreedsData handles concurrent data loading for breeds view
func (c *UnifiedController) GetBreedsData() {
    var wg sync.WaitGroup
    var mutex sync.Mutex
    data := BreedsViewData{}
    errors := []string{}
    
    breedsChan := make(chan []Breed, 1)
    imagesChan := make(chan []Image, 1)
    
    // Fetch breeds
    wg.Add(1)
    go func() {
        defer wg.Done()
        breeds, err := fetchBreeds()
        if err != nil {
            mutex.Lock()
            errors = append(errors, "Failed to fetch breeds: "+err.Error())
            mutex.Unlock()
            return
        }
        breedsChan <- breeds
    }()
    
    // Fetch breed images if breed_id is provided
    breedID := c.GetString("breed_id")
    if breedID != "" {
        wg.Add(1)
        go func() {
            defer wg.Done()
            images, err := fetchBreedImages(breedID)
            if err != nil {
                mutex.Lock()
                errors = append(errors, "Failed to fetch breed images: "+err.Error())
                mutex.Unlock()
                return
            }
            imagesChan <- images
        }()
    }
    
    // Handle completion or timeout
    done := make(chan bool)
    go func() {
        wg.Wait()
        done <- true
    }()
    
    select {
    case <-done:
        select {
        case data.Breeds = <-breedsChan:
        default:
        }
        select {
        case data.Images = <-imagesChan:
        default:
        }
    case <-time.After(12 * time.Second):
        errors = append(errors, "Operation timed out")
    }
    
    if len(errors) > 0 {
        data.Errors = errors
    }
    
    c.Data["json"] = data
    c.ServeJSON()
}
func (c *UnifiedController) fetchFavorites(apiKey, subID string) ([]Favorite, error) {
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

func (c *UnifiedController) fetchImages(apiKey string) ([]models.CatImage, error) {
    client := &http.Client{Timeout: 10 * time.Second}
    req, _ := http.NewRequest("GET", 
        fmt.Sprintf("https://api.thecatapi.com/v1/images/search?limit=10"), 
        nil)
    
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