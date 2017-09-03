package models

import (
	"regexp"
	"time"

	"github.com/revel/revel"
)

const (
	//EventTitleMaxLength is the maximum allowed length for title field
	EventTitleMaxLength = 500

	//EventAddressMaxLength is the maximum allowed length for address field
	EventAddressMaxLength = 500

	//EventCategoryMaxLength is the maximum allowed length for category field
	EventCategoryMaxLength = 50

	//EventMessageMaxLength is the maximum allowed length for message field
	EventMessageMaxLength = 1000

	//EventLinkleMaxLength is the maximum allowed length for link field
	EventLinkleMaxLength = 500

	//EventTagMaxLength is the maximum allowed length for tag field
	EventTagMaxLength = 50
)

// Event model
type Event struct {
	ID          int64     `db:"id" json:"id"`
	Tag         string    `db:"tag" json:"tag"`
	Title       string    `db:"title" json:"title"`
	Category    string    `db:"category" json:"category"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	Message     string    `db:"message" json:"message"`
	Address     string    `db:"address" json:"address"`
	Longitude   float32   `db:"longitude" json:"longitude"`
	Latitude    float32   `db:"latitude" json:"latitude"`
	MapLink     string    `db:"map_link" json:"mapLink"`
	Time        time.Time `db:"time" json:"time"`
	MinAttendee int16     `db:"min_attendee" json:"minAttendee"`
	MaxAttendee int16     `db:"max_attendee" json:"maxAttendee"`
	Status      string    `db:"status" json:"status"`
}

// Validate method validates the Event fields
func (e *Event) Validate(v *revel.Validation) {

	v.Check(e.Title,
		revel.ValidRequired(),
		revel.ValidMaxSize(EventTitleMaxLength))

	v.Check(e.Category,
		revel.ValidRequired(),
		revel.ValidMatch(
			regexp.MustCompile(
				"^(futsal|other)$")))

	v.Check(e.Status, revel.ValidRequired(),
		revel.ValidMatch(
			regexp.MustCompile(
				"^(active|cancelled|deleted)$")))

	v.Check(e.Time,
		revel.ValidRequired())
}
