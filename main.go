package main

import (
	"net/http"

	"github.com/ccclin/gae_esun/controller"
	"google.golang.org/appengine/v2"
)

func main() {
	http.HandleFunc("/check", controller.CheckHandle)
	http.HandleFunc("/queue", controller.QueueHandle)

	http.HandleFunc("/send", controller.SendHandle)
	appengine.Main()
}
