package main

import "lsq.com/router"

func main() {
	r := router.Router()

	r.Run(":8080")
}
