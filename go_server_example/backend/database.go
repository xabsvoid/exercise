package main

import (
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	con *gorm.DB
}

func (db *Database) Init(dsn string) error {
	var err error
	db.con, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	err = db.con.AutoMigrate(&Note{})
	if err != nil {
		return err
	}
	return nil
}

func (db Database) SelectNotes(offset uint64, limit uint16) ([]Note, error) {
	var notes []Note
	res := db.con.Where("id >= ?", offset).Limit(int(limit)).Find(&notes)
	return notes, res.Error
}

func (db Database) SelectNote(id uint64) (Note, error) {
	var note Note
	var notes []Note
	res := db.con.Where("id = ?", id).Find(&notes)
	if res.Error != nil {
		return note, res.Error
	}
	if len(notes) == 0 {
		return note, errors.New("note not found")
	}
	return notes[0], nil
}

func (db Database) InsertNote(note Note) error {
	res := db.con.Create(&note)
	return res.Error
}

func (db Database) UpdateNote(note Note) error {
	res := db.con.Model(&note).Where("id = ?", note.ID).Updates(&note)
	return res.Error
}

func (db Database) DeleteNote(id uint64) error {
	var note Note
	res := db.con.Delete(&note, id)
	return res.Error
}
