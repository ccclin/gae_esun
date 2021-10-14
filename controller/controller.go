package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ccclin/gae_esun/model"
	"google.golang.org/appengine/v2"
	"google.golang.org/appengine/v2/log"
	"google.golang.org/appengine/v2/taskqueue"
)

// CheckHandle is GET '/check'
func CheckHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

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

// QueueHandle is POST '/queue'
func QueueHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/queue" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	ctx := appengine.NewContext(r)
	mail := model.Mail{Ctx: ctx}
	mail.GetEsun()
	mail.Send()
}

// SendHandle is POST '/send'
func SendHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/send" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	ctx := appengine.NewContext(r)
	var expected model.Expected
	err := json.NewDecoder(r.Body).Decode(&expected)
	if err != nil {
		log.Errorf(ctx, "son.NewDecoder(r.Body).Decode failed %+v", err)
		errorHandler(w, r, http.StatusBadRequest)
		return
	}
	log.Infof(ctx, "expected is %f", expected.Expected)
	expected.PutDatastore(ctx, &expected)
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
