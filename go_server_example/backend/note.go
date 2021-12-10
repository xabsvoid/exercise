package main

import (
	"errors"
	"time"
)

type Note struct {
	ID          uint64    `json:"id" gorm:"primary_key"`
	CreatedTime time.Time `json:"created_time"`
	Title       string    `json:"title"`
	Text        string    `json:"text"`
}

func (Note) TableName() string {
	return "note"
}

func (note Note) Validate() error {
	if len(note.Title) > 128 {
		return errors.New("title too long")
	}
	if len(note.Text) > 4096 {
		return errors.New("text too long")
	}
	return nil
}
