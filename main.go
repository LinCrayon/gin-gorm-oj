package main

import "github.com/LinCrayon/gin-gorm-oj/router"

func main() {
	r := router.Router()

	r.Run(":8080")
}
