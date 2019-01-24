package model

import (
	"context"
	"os"
	"strconv"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/mail"
	"google.golang.org/appengine/memcache"
)

// SENDKEY is for queue to send mail
const SENDKEY = "SEND_KEY"

// Mail for send mail
type Mail struct {
	Ctx  context.Context
	Esun Esun
}

// GetEsun is get ESUN from memcache
func (m *Mail) GetEsun() {
	memcache.JSON.Get(m.Ctx, SENDKEY, &m.Esun)
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
	mail.Send(ctx, msg)
}
