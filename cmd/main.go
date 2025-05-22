package main

import (
	"fmt"
	router "github.com/niuzheng1314520/gin/api/routes"
)

func main() {
	r := router.Router()
	err := r.Run(":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
}
