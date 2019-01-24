package controller

import (
	"fmt"
	"net/http"

	"github.com/ccclin/gae_esun/model"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
)

// CheckHandle is GET '/check'
func CheckHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/check" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	ctx := appengine.NewContext(r)
	esun := model.Esun{Ctx: ctx}
	c := make(chan bool)
	go esun.SetExpected(c)
	go esun.GetJPY(c)
	<-c
	<-c

	log.Infof(ctx, "JPY is %v", esun.JPY)
	log.Infof(ctx, "Expected is %v", esun.Expected)

	if esun.JPY < esun.Expected {
		esun.SetMemcache()
		t := taskqueue.NewPOSTTask("/queue", nil)
		if _, err := taskqueue.Add(ctx, t, "send-mail"); err != nil {
			errorHandler(w, r, http.StatusInternalServerError)
			return
		}
	}

	if esun.Err != nil {
		log.Errorf(ctx, "err is %s", esun.Err)
		errorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}
}

// QueueHandle is POST '/post'
func QueueHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/queue" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	ctx := appengine.NewContext(r)
	mail := model.Mail{Ctx: ctx}
	mail.GetEsun()
	mail.Send()
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	switch status {
	case http.StatusNotFound:
		fmt.Fprint(w, "404 Not Found")
	case http.StatusMethodNotAllowed:
		fmt.Fprint(w, "405 Method Not Allow")
	default:
		fmt.Fprint(w, "Bad Request")
	}
}
