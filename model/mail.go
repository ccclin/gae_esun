package model

import (
	"context"
	"os"
	"strconv"
	"strings"

	"google.golang.org/appengine/v2"
	"google.golang.org/appengine/v2/log"
	"google.golang.org/appengine/v2/mail"
	"google.golang.org/appengine/v2/memcache"
)

// SendKey is for queue to send mail
const SendKey = "SEND_KEY"

// Mail for send mail
type Mail struct {
	Ctx  context.Context
	Esun Esun
}

// GetEsun is get ESUN from memcache
func (m *Mail) GetEsun() {
	_, err := memcache.JSON.Get(m.Ctx, SendKey, &m.Esun)
	if err != nil {
		log.Errorf(m.Ctx, "memcache.JSON.Get failed %+v", err)
	}
}

// Send will send mail
func (m *Mail) Send() {
	ctx := m.Ctx
	jpyPricing := strconv.FormatFloat(m.Esun.JPY, 'f', -1, 64)
	htmlBody := `JPY is ` + jpyPricing

	msg := &mail.Message{
		Sender:   "noreply@" + appengine.AppID(ctx) + ".appspotmail.com",
		To:       strings.Split(os.Getenv("TO"), ","),
		Subject:  "[ESUN] JPY is " + jpyPricing,
		HTMLBody: htmlBody,
	}
	err := mail.Send(ctx, msg)
	if err != nil {
		log.Errorf(ctx, "send mail failed %+v", err)
	}
}
