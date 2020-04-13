package main

import (
	"github.com/ontio/sagapi/restful"
)

func main() {
	router := restful.NewRouter()
	router.Run(":8080")
}
