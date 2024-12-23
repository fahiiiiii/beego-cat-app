package main

import (
	_ "beego-cat-app/routers" // Import your routers package
	"github.com/beego/beego/v2/server/web"
)

func main() {
	web.Run() // Start the Beego server
}
