package model

import (
	"context"
	"encoding/binary"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/urlfetch"
)

// ESUNJPY is esun jpy now prcing
const (
	ESUNJPY = "ESUN_JPY"
)

// Esun is esun bank
type Esun struct {
	Expected float64
	JPY      float64
	Err      error
	Ctx      context.Context
}

// SetExpected is set expected to memcache
func (esun *Esun) SetExpected(c chan bool) {
	expJPY, _ := strconv.ParseFloat(os.Getenv("EXPECTED"), 64)
	esun.getFromMemcache(expJPY)
	c <- true
}

// GetJPY is get JPY from ESUN
func (esun *Esun) GetJPY(c chan bool) {
	ctx := esun.Ctx
	client := urlfetch.Client(ctx)
	req, _ := http.NewRequest("GET", os.Getenv("ESUN"), nil)

	resp, err := client.Do(req)
	if err != nil {
		esun.setErrorAndJPY(err, 1.0)
		c <- true
		return
	}

	crawler := Crawler{Resp: resp}
	jpy, err := crawler.GetJpy()
	if err != nil {
		esun.setErrorAndJPY(err, 1.0)
		c <- true
		return
	}
	esun.JPY = jpy
	c <- true
}

func (esun *Esun) getFromMemcache(oriJPY float64) {
	value, err := memcache.Get(esun.Ctx, ESUNJPY)
	if err == nil {
		esun.Expected = ByteToFloat64(value.Value)
	} else {
		esun.Expected = oriJPY
	}
}

func (esun *Esun) setErrorAndJPY(err error, jpy float64) {
	log.Errorf(esun.Ctx, "err %s", err)
	esun.Err = err
	esun.JPY = jpy
}

// SetMemcache is set into memcache
func (esun *Esun) SetMemcache() {
	esun.setESUNJPY()
	esun.setMail()
}

func (esun *Esun) setESUNJPY() {
	item := &memcache.Item{
		Key:        ESUNJPY,
		Value:      Float64ToByte(esun.JPY),
		Expiration: time.Duration(28800) * time.Second,
	}

	if err := memcache.Set(esun.Ctx, item); err != nil {
		esun.Err = err
	}
}

func (esun *Esun) setMail() {
	item := &memcache.Item{
		Key:    SENDKEY,
		Object: esun,
	}

	if err := memcache.JSON.Set(esun.Ctx, item); err != nil {
		esun.Err = err
	}
}

// Float64ToByte is float64 to byte
func Float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)

	return bytes
}

// ByteToFloat64 is byte to float64
func ByteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)

	return math.Float64frombits(bits)
}
