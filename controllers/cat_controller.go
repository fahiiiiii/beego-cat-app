//controllers/cat_controller.go

package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/beego/beego/v2/server/web"
)

type CatController struct {
	web.Controller
}

type Cat struct {
	ID     string `json:"id"`
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// FetchCats fetches cat images from The Cat API
func (c *CatController) FetchCats() {
    apiKey, err := web.AppConfig.String("cat_api_key")
    if err != nil || apiKey == "" {
        fmt.Println("API key is not set in app.conf")
        c.Data["json"] = map[string]string{"error": "API key not found"}
        c.ServeJSON()
        return
    }

    baseURL, err := web.AppConfig.String("https://api.thecatapi.com/v1")
    if err != nil || baseURL == "" {
        fmt.Println("Base URL is not set in app.conf")
        c.Data["json"] = map[string]string{"error": "Base URL not found"}
        c.ServeJSON()
        return
    }

    url := fmt.Sprintf("%s/images/search", baseURL)

    dataChan := make(chan []Cat)

    go func() {
        client := &http.Client{}
        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
            fmt.Println("Error creating request:", err)
            dataChan <- nil
            return
        }
        req.Header.Set("x-api-key", apiKey)

        resp, err := client.Do(req)
        if err != nil {
            fmt.Println("Error fetching cats:", err)
            dataChan <- nil
            return
        }
        defer resp.Body.Close()

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            fmt.Println("Error reading response:", err)
            dataChan <- nil
            return
        }

        // Debug: Log the response body
        fmt.Println("Response body:", string(body))

        var cats []Cat
        if err := json.Unmarshal(body, &cats); err != nil {
            fmt.Println("Error unmarshalling JSON:", err)
            dataChan <- nil
            return
        }
        dataChan <- cats
    }()

    cats := <-dataChan
    if cats == nil || len(cats) == 0 {
        c.Data["json"] = map[string]string{"error": "No cats found"}
    } else {
        c.Data["json"] = cats
    }
    c.ServeJSON()
}
