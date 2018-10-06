package app

import (
	"context"
	"os"
	"strconv"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/mail"
	"google.golang.org/appengine/memcache"
)

// Mail for send mail
type Mail struct {
	Ctx  context.Context
	Esun Esun
}

func (m *Mail) getEsun() {
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
