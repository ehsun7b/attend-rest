package controllers

import (
	"encoding/json"
	"time"

	"github.com/ehsun7b/attend-rest/app/models"
	"github.com/revel/revel"
	hashids "github.com/speps/go-hashids"
)

var (
	newTagsLength = 8
)

// EventCtrl the controller of Event
type EventCtrl struct {
	GorpController
}

func (c EventCtrl) parseEvent() (models.Event, error) {
	event := models.Event{}
	err := json.NewDecoder(c.Request.Body).Decode(&event)
	return event, err
}

//Add inserts new event object
func (c EventCtrl) Add() revel.Result {
	event, err := c.parseEvent()

	event.CreatedAt = time.Now()

	if err != nil {
		// todo log
		msg := "Unable to parse the Event from JSON."
		revel.WARN.Printf(msg)
		c.Response.Status = 400
		return c.RenderText(msg)
	}
	// Validate the model
	event.Validate(c.Validation)
	if c.Validation.HasErrors() {
		msg := "You have error in your Event."
		c.Response.Status = 400
		revel.WARN.Printf(msg)
		return c.RenderText(msg)
	}

	if err := c.Txn.Insert(&event); err != nil {
		c.Response.Status = 500
		return c.RenderText("Error inserting record into database!" + err.Error())
	}

	event.Tag = eventTag(event)
	_, err = c.Txn.Update(&event)
	tries := 0

	for err != nil && tries < 10 {
		event.Tag = eventTag(event)
		revel.TRACE.Println(event.Tag)
		_, err = c.Txn.Update(&event)
		tries++
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		c.Txn.Delete(&event)
		c.Response.Status = 500
		return c.RenderText("Error inserting record into database!" + err.Error())
	}

	return c.RenderJSON(event)
}

func eventTag(event models.Event) string {
	hd := hashids.NewData()
	salt, found := revel.Config.String("attend.hashid.salt")

	if !found {
		revel.ERROR.Panic("attend.hashid.salt not found in the app.conf")
	}

	len := newTagsLength

	confLen, found2 := revel.Config.Int("attend.hashid.length")

	if !found2 {
		revel.WARN.Println("attend.hashid.length not found in the app.conf")
	} else {
		len = confLen
	}

	hd.Salt = salt
	hd.MinLength = len
	h, _ := hashids.NewWithData(hd)
	e, _ := h.EncodeInt64([]int64{event.ID})
	return e
}
