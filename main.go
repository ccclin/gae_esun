package main

import (
	"net/http"

	"github.com/ccclin/gae_esun/controller"
	"google.golang.org/appengine"
)

func main() {
	http.HandleFunc("/check", controller.CheckHandle)
	http.HandleFunc("/queue", controller.QueueHandle)
	appengine.Main()
}
