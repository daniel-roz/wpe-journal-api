package models

import (
	"uuid"
)

type Category string

const (
	CHURCH    Category = "church"
	CLASSROOM Category = "classroom"
	COURT     Category = "court"
	COMMUNITY Category = "community"
	CULTURE   Category = "culture"
)

func (c Category) IsValid() bool {
	switch c {
	case
		CHURCH,
		CLASSROOM,
		COURT,
		COMMUNITY,
		CULTURE:
		return true
	default:
		return false
	}
}

type Journal struct {
	Uuid         uuid.UUID `json:"uuid"`
	Title        string    `json:"title"`
	Category     Category  `json:"category"`
	Organization string    `json:"organization"`
	Date         string    `json:"date"`
	Venue        string    `json:"venue"`
	Role         string    `json:"role"`
	Points       int       `json:"points"`
	Article      string    `json:"article"`
	Images       []string  `json:"images"`
}
type JournalUpdate struct {
	Title        *string   `json:"title"`
	Category     *Category `json:"category"`
	Organization *string   `json:"organization"`
	Date         *string   `json:"date"`
	Venue        *string   `json:"venue"`
	Role         *string   `json:"role"`
	Points       *int      `json:"points"`
	Article      *string   `json:"article"`
	Images       *[]string `json:"images"`
}
type JournalPreview struct {
	Uuid         uuid.UUID `json:"uuid"`
	Title        string    `json:"title"`
	Organization string    `json:"organization"`
}
