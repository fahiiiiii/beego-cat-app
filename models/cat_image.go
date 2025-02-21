// models/cat_image.go
package models

// CatImage represents the structure of a cat image from the CatAPI
type CatImage struct {
    ID  string `json:"id"`
    URL string `json:"url"`
}
