package app

import (
	"net/http"

	"google.golang.org/appengine"
)

// SENDKEY is for queue to send mail
const SENDKEY = "SEND_KEY"

func init() {
	http.HandleFunc("/check", checkHandle)
	http.HandleFunc("/queue", queueHandle)
	appengine.Main()
}
