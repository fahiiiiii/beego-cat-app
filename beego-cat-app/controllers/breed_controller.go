// controllers/breed_controller.go

package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/server/web"
	// "log"
	"net/http"
)

type BreedController struct {
	web.Controller
}

type Breed struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Origin      string `json:"origin"`
	WikipediaURL string `json:"wikipedia_url"`
}

type Image struct {
	URL string `json:"url"`
}

var breedAPIURL = "https://api.thecatapi.com/v1/breeds"
var imageAPIURL = "https://api.thecatapi.com/v1/images/search"

// Fetch breeds from the external API
func fetchBreeds() ([]Breed, error) {
	resp, err := http.Get(breedAPIURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var breeds []Breed
	if err := json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
		return nil, err
	}
	return breeds, nil
}

// Fetch breed images from the external API
func fetchBreedImages(breedID string) ([]Image, error) {
	url := fmt.Sprintf("%s?breed_ids=%s&limit=10", imageAPIURL, breedID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var images []Image
	if err := json.NewDecoder(resp.Body).Decode(&images); err != nil {
		return nil, err
	}
	return images, nil
}

// Action to get breeds (mapped to /api/breeds)
func (c *BreedController) GetBreeds() {
	breeds, err := fetchBreeds()
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Ctx.WriteString("Failed to fetch breeds")
		return
	}

	c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(c.Ctx.ResponseWriter).Encode(breeds); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Ctx.WriteString("Failed to encode breeds data")
	}
}

// Action to get breed images (mapped to /api/breed-images)
func (c *BreedController) GetBreedImages() {
	breedID := c.GetString("breed_id")
	if breedID == "" {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Ctx.WriteString("Breed ID is required")
		return
	}

	images, err := fetchBreedImages(breedID)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Ctx.WriteString("Failed to fetch breed images")
		return
	}

	c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(c.Ctx.ResponseWriter).Encode(images); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Ctx.WriteString("Failed to encode breed images")
	}
}
